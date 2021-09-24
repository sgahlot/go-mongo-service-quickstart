.DEFAULT_GOAL := help

GO_CMD := go
GO_RUN := go run cmd/main.go
BINARY_NAME := mongo-svc-quickstart

# Overwrite this variable from cli to specify private different OS
TARGET_OS ?= linux

QUAY_USER_OR_ORG ?= sgahlot
IMAGE_NAME_WITHOUT_TAG ?= go-mongo-quickstart
IMAGE_TAG ?= 0.0.1-SNAPSHOT
DB_URL ?= ""
DB_NAME ?= fruit
SERVICE_BINDING_ROOT ?= ${PWD}/test-bindings/bindings

IMAGE_REGISTRY ?= quay.io
IMAGE_REPO ?= $(USER)
IMAGE_NAME ?= $(IMAGE_NAME_WITHOUT_TAG):$(IMAGE_TAG)
CONTAINER_NAME ?= go-mongo-fruit-app
IMAGE_BUILDER ?= docker

## Include and export the environment variables
#include resources/docker/go/.env
#export

DOCKER_ENV = QUAY_USER_OR_ORG=$(QUAY_USER_OR_ORG) \
              IMAGE_NAME_WITHOUT_TAG=$(IMAGE_NAME_WITHOUT_TAG) \
              IMAGE_TAG=$(IMAGE_TAG) \
              DB_URL=$(DB_URL) \
              DB_NAME=$(DB_NAME) \
              CONTAINER_NAME=$(CONTAINER_NAME) \
              BINARY_NAME=$(BINARY_NAME) \
              SERVICE_BINDING_ROOT=$(SERVICE_BINDING_ROOT)

.PHONY: run
run: ## Runs the app from source - without building an executable
	SERVICE_BINDING_ROOT=$(SERVICE_BINDING_ROOT) $(GO_RUN)

.PHONY: build_binary
build_binary: delete_binary ## Creates the binary for the app after deleting previous one as it could be for different OS
	CGO_ENABLE=0 GOOS=$(TARGET_OS) GOARCH=amd64 go build -o $(BINARY_NAME) ./cmd

.PHONY: delete_binary
delete_binary: ## Creates the binary for the app
	rm -f $(BINARY_NAME)

.PHONY: run_binary
run_binary: build_binary ## Creates the binary for the app and runs it
	SERVICE_BINDING_ROOT=$(SERVICE_BINDING_ROOT) ./$(BINARY_NAME)

.PHONY: show_build_config
show_build_config: ## Shows the docker compose output - with all the environment variable substitution
	$(DOCKER_ENV) $(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml config

.PHONY: build_image
build_image: build_binary ## Builds docker image
	$(DOCKER_ENV) $(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml build

.PHONY: start_container
start_container: build_image stop_mongo_db ## Starts the Docker container (and builds docker image if not already built)
	$(DOCKER_ENV) $(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml up -d

.PHONY: stop_container
stop_container: ## Stops the docker container and removes it as well
	$(DOCKER_ENV) $(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml down

.PHONY: start_mongo_db
start_mongo_db: ## Runs ONLY MongoDB container
	$(DOCKER_ENV) $(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml up -d mongo_db

.PHONY: stop_mongo_db
stop_mongo_db: ## Stops MongoDB container
	$(DOCKER_ENV) $(IMAGE_BUILDER) compose -f go-svc-docker-compose.yaml down mongo_db

.PHONY: push_image
push_image: ## Builds the image and pushes it to quay.io
	$(IMAGE_BUILDER) push $(IMAGE_REGISTRY)/$(IMAGE_REPO)/$(IMAGE_NAME)

