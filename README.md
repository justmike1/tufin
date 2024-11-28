# Deployer CLI Tool

## Prerequisites

1. Install Dependencies
    ```bash
    go mod download
    ```

2. Build the CLI tool
    ```bash
    go build cmd/main/deployer.go
    ```
   
## Run the CLI tool

**Deploy the k3s cluster**

```bash
./deployer cluster
```

**Deploy the application**

```bash
./deployer deploy
```

**Check the status of the application**

```bash
./deployer status
```

**Open the application**

http://localhost:30080/


