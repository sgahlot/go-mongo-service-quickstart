.DEFAULT_GOAL := help

GO_CMD := go
GO_RUN := go run cmd/main.go
EXE_NAME := mongo-svc-quickstart

# Overwrite this variable from cli to specify private different OS
EXE_OS ?= darwin

IMAGE_REGISTRY ?= quay.io
IMAGE_REPO ?= $(USER)
IMAGE_NAME ?= go-mongo-quickstart:0.0.1-SNAPSHOT
IMAGE_BUILDER ?= docker

## Include and export the environment variables
include resources/docker/go/.env
export

.PHONY: run
run: ## Runs the app from source - without building an executable
	$(GO_RUN)

.PHONY: run_binary
run_binary: ## Creates the binary for the app and runs it
	CGO_ENABLE=0 GOOS=$(EXE_OS) GOARCH=amd64 go build -o $(EXE_NAME) ./cmd
	./$(EXE_NAME)

.PHONY: build_config
build_config: ## Shows the docker compose output - with all the environment variable substitution
	$(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml config

.PHONY: build_image
build_image: ## Builds docker image
	$(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml build

.PHONY: run_app
run_app: build_image ## Builds docker image as well as starts the container
	$(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml up -d

.PHONY: stop_app
stop_app: ## Builds docker image as well as the container
	$(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml down

.PHONY: push_image
push_image: ## Builds the image and pushes it to quay.io
	$(IMAGE_BUILDER) push $(IMAGE_REGISTRY)/$(IMAGE_REPO)/$(IMAGE_NAME)

