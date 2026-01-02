TAG=0.0.2

.PHONY: build
build: build

build:
	docker buildx build --platform "linux/arm64,linux/amd64" -t ep4sh/pastecode:$(TAG) --push .

all: build
