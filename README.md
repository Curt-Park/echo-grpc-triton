# echo-grpc-triton
A simple API server for Triton inference server powered by echo and gRPC.

## Prerequisites
- Install [docekr](https://docs.docker.com/engine/install/).
- Install [go](https://go.dev/doc/install).

## Commands
```bash
make triton-run             # Run triton server
docker logs tritonserver    # See Triton server's logs
go run *.go                 # Run the API server
make triton-kill            # Kill triton server
```

Open http://localhost:8080/docs/index.html to see the API documents.

## References
- https://github.com/triton-inference-server/client/tree/main/src/grpc_generated/go
- https://github.com/sunhailin-Leo/triton-service-go
- https://echo.labstack.com/guide/
- https://github.com/swaggo/swag#api-operation
