TAG=0.0.1

.PHONY: build
build: build

build:
	docker buildx build --platform "linux/arm64,linux/amd64" -t ep4sh/pastecode:$(TAG) --push .

all: build
