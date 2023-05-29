# Run a local development Temporal Server

https://docs.temporal.io/dev-guide/go/foundations#run-a-development-server

## macOS

    brew install temporal
    temporal server start-dev
    
## Linux

    curl -sSf https://temporal.download/cli.sh | sh
    temporal server start-dev
    
# Run Expense example

* Start the dummy server 
```
go run expense/server/main.go
```
* Start workflow and activity workers
```
go run expense/worker/main.go
```
* Start expanse workflow execution
```
go run expense/starter/main.go
```    