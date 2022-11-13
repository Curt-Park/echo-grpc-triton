# echo-grpc-triton
A simple API server for Triton inference server powered by echo and gRPC.

## Prerequisites
- Install [docker](https://docs.docker.com/engine/install/).
- Install [go](https://go.dev/doc/install) (Optional: If you want to run this project without Docker).

## Commands
```bash
docker-compose up           # Run all services (API + Triton)
docker-compose kill         # Kill all services
```

Open http://localhost:8080/docs/index.html to see the API documents and to send a single request.

<img width="1352" src="https://user-images.githubusercontent.com/14961526/201515996-0e03f353-4017-42aa-b704-d8ac140a6e0f.png">

Press `Try it out` button to call any API.
<img width="1350" src="https://user-images.githubusercontent.com/14961526/201516009-e18e966f-ae8b-4bcd-8f21-429c3152301d.png">


## References
- https://github.com/triton-inference-server/client/tree/main/src/grpc_generated/go
- https://github.com/sunhailin-Leo/triton-service-go
- https://echo.labstack.com/guide/
- https://github.com/swaggo/swag#api-operation
