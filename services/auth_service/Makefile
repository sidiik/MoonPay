# Variables
DOCKER_BUILD_CONTEXT = .
IMAGE_PREFIX = moonpay
SERVICE = auth-service
VERSION = 0.0.1
IMAGE = $(IMAGE_PREFIX)/$(SERVICE):$(VERSION)

docker-build:
	docker build \
		--build-arg SERVICE=$(SERVICE) \
		--build-arg BINARY_NAME=$(SERVICE) \
		-t $(IMAGE) \
		$(DOCKER_BUILD_CONTEXT)

docker-push:
	docker push $(IMAGE)

docker-release: docker-build docker-push

protoc:
	protoc --go_out=internal/delivery/grpc --go_opt=paths=source_relative \
    --go-grpc_out=internal/delivery/grpc --go-grpc_opt=paths=source_relative \
    proto/auth.proto
