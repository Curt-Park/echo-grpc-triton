triton-run:
	docker run -d --name tritonserver -p8000:8000 -p8001:8001 -p8002:8002 -it -v $(PWD)/models:/models nvcr.io/nvidia/tritonserver:22.09-py3 tritonserver --model-store=/models

triton-kill:
	docker stop tritonserver
	docker rm tritonserver

format:
	go fmt
	swag fmt
	swag init
