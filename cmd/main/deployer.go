package main

import (
	"github.com/justmike1/deployer/pkg/cluster"
	"github.com/justmike1/deployer/pkg/config"
	"github.com/justmike1/deployer/pkg/deploy"
	"github.com/justmike1/deployer/pkg/status"
	"log"
	"os"
)

const (
	clusterName = "test"
	namespace   = "default"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("No command provided. Use 'help' command to see available commands.")
		os.Exit(1)
	}

	cmdOption := config.ParseCmdOption(os.Args[1])

	switch cmdOption {
	case config.CLUSTER:
		deployK3s()
	case config.DEPLOY:
		deployK8sManifests()
	case config.STATUS:
		showStatus()
	default:
		log.Println("Unknown command. Use 'help' command to see available commands.")
		os.Exit(1)
	}
}

func deployK3s() {
	log.Println("Deploying a local K3s cluster...")
	cluster.K3dCluster(clusterName)
}

func deployK8sManifests() {
	log.Println("Deploying a local WordPress and MySQL application manifests...")
	deploy.WordPressAndMySQL(clusterName, namespace)
}

func showStatus() {
	log.Println("Checking the application status...")
	status.LogPodStatuses(clusterName, namespace)
}
