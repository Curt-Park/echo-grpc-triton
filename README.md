# echo-grpc-triton
A simple API server for Triton inference server powered by echo and gRPC.

## Prerequisites
- Install [go](https://go.dev/doc/install).

## Commands
```bash
make triton-run             # Run triton server
make triton-kill            # Kill triton server
docker logs tritonserver    # See Triton server's logs
go run *.go              # Run the API server
```

Open http://localhost:8080/docs/index.html to see the API documents.
