.PHONY: build test install-dev docker api

PROTO_BUILD_DIR = ./pkg

#TEST_ARGS = -v

DOCKER_OPTS ?= --rm
VERSION := $(shell git describe --tags --abbrev=0)

help:
	@echo "Service building targets"
	@echo "  build : build service command"
	@echo "  test  : run test suites"
	@echo "  docker: build docker image"
	@echo "  install-dev: install dev dependencies"
	@echo "  api: compile protobuf files for go"
	@echo "Env:"
	@echo "  DOCKER_OPTS : default docker build options (default : $(DOCKER_OPTS))"
	@echo "  TEST_ARGS : Arguments to pass to go test call"

api:
	if [ ! -d "$(PROTO_BUILD_DIR)/api" ]; then mkdir -p "$(PROTO_BUILD_DIR)"; else  find "$(PROTO_BUILD_DIR)/api" -type f -delete &&  mkdir -p "$(PROTO_BUILD_DIR)"; fi
	find ./api/*.proto -maxdepth 1 -type f -exec protoc {} --go_opt=paths=source_relative --go_out=plugins=grpc:$(PROTO_BUILD_DIR) \;

build:
	go build .

test:
	./test/test.sh $(TEST_ARGS)

docker:
	docker build -t github.com/influenzanet/study-service:$(VERSION)  -f build/docker/Dockerfile $(DOCKER_OPTS) .