# Tufin CLI Tool

## Prerequisites

1. Install Dependencies
    ```bash
    go mod download
    ```

2. Build the CLI tool
    ```bash
    go build cmd/main/tufin.go
    ```
   
## Run the CLI tool

**Deploy the k3s cluster**

```bash
./tufin cluster
```

**Deploy the application**

```bash
./tufin deploy
```

**Check the status of the application**

```bash
./tufin status
```

**Open the application**

http://localhost:30080/


