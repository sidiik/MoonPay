# Variables
DOCKER_BUILD_CONTEXT = .
IMAGE_PREFIX = moonpay
SERVICE = api-gateway
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
