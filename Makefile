ORG ?= crazytaxii
IMG_NAME ?= ${ORG}/fake-dcgm-exporter
IMG_TAG ?= latest
TARGET_DIR ?= .dist

build:
	CGO_ENABLED=0 go build -o $(TARGET_DIR)/fake-dcgm-exporter ./cmd/

image:
	docker buildx build -t $(IMG_NAME):$(IMG_TAG) --push --platform linux/amd64,linux/arm64 .

PHONY: build image
