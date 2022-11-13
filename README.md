# echo-grpc-triton
A simple API server for Triton inference server powered by echo and gRPC.

## Prerequisites
- Install [docker](https://docs.docker.com/engine/install/).
- Install [go](https://go.dev/doc/install).

## Commands
```bash
docker-compose up           # Run all services (API + Triton)
docker-compose kill         # Kill all services
```

Open http://localhost:8080/docs/index.html to see the API documents.

## References
- https://github.com/triton-inference-server/client/tree/main/src/grpc_generated/go
- https://github.com/sunhailin-Leo/triton-service-go
- https://echo.labstack.com/guide/
- https://github.com/swaggo/swag#api-operation
