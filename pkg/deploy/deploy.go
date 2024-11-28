package deploy

import (
	"context"
	"github.com/justmike1/deployer/pkg/cluster"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func WordPressAndMySQL(clusterName string, namespace string) {
	kubeconfigContent, err := cluster.GetKubeconfigContent(clusterName)
	if err != nil {
		log.Fatalf("Failed to get kubeconfig for cluster %s: %v", clusterName, err)
	}

	kubeconfigPath, err := cluster.CreateTempKubeconfigFile(kubeconfigContent)
	if err != nil {
		log.Fatalf("Failed to create temporary kubeconfig file: %v", err)
	}
	defer os.Remove(kubeconfigPath)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes clientset: %v", err)
	}

	if err := createMySQLDeployment(clientset, namespace); err != nil {
		log.Fatalf("Failed to create MySQL deployment: %v", err)
	}
	log.Println("MySQL deployment and service created.")

	if err := createWordPressDeployment(clientset, namespace); err != nil {
		log.Fatalf("Failed to create WordPress deployment: %v", err)
	}
	log.Println("WordPress deployment and service created.")
}

func createMySQLDeployment(clientset *kubernetes.Clientset, namespace string) error {
	// MySQL Deployment
	mysqlDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mysql",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "mysql"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "mysql"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "mysql",
							Image: "mysql:lts",
							Env: []corev1.EnvVar{
								{Name: "MYSQL_ROOT_PASSWORD", Value: "password"},
								{Name: "MYSQL_DATABASE", Value: "wordpress"},
								{Name: "MYSQL_USER", Value: "wordpress"},
								{Name: "MYSQL_PASSWORD", Value: "wordpress"},
							},
						},
					},
				},
			},
		},
	}
	if _, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), mysqlDeployment, metav1.CreateOptions{}); err != nil {
		return err
	}

	// MySQL Service
	mysqlService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mysql",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "mysql"},
			Ports: []corev1.ServicePort{
				{Port: 3306, TargetPort: intstr.FromInt(3306)}, // Updated here
			},
		},
	}
	if _, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), mysqlService, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func createWordPressDeployment(clientset *kubernetes.Clientset, namespace string) error {
	// WordPress Deployment
	wordpressDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "wordpress",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "wordpress"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "wordpress"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "wordpress",
							Image: "wordpress:php7.4-apache",
							Env: []corev1.EnvVar{
								{Name: "WORDPRESS_DB_HOST", Value: "mysql:3306"},
								{Name: "WORDPRESS_DB_USER", Value: "wordpress"},
								{Name: "WORDPRESS_DB_PASSWORD", Value: "wordpress"},
								{Name: "WORDPRESS_DB_NAME", Value: "wordpress"},
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 80},
							},
						},
					},
				},
			},
		},
	}
	if _, err := clientset.AppsV1().Deployments(namespace).Create(context.TODO(), wordpressDeployment, metav1.CreateOptions{}); err != nil {
		return err
	}

	// WordPress Service
	wordpressService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "wordpress",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "wordpress"},
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt32(80),
					NodePort:   30080, // Chosen NodePort (use any port in range 30000-32767)
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}
	if _, err := clientset.CoreV1().Services(namespace).Create(context.TODO(), wordpressService, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
