# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOTIDY = $(GOCMD) mod tidy
BINARY_NAME = social_network
MAIN = cmd/server/main.go

# Docker parameters
DOCKER = docker
DOCKER_BUILD = $(DOCKER) build
DOCKER_RUN = $(DOCKER) run
DOCKER_IMAGE = social_network:latest

# Commands
all: test build

build:
	$(GOTIDY)
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN)

test:
	$(GOTEST) -v ./...

run:
	$(GOTIDY)
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN)
	./$(BINARY_NAME)

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)

docker-build:
	$(DOCKER_BUILD) -t $(DOCKER_IMAGE) .

docker-run:
	$(DOCKER_RUN) -p 8080:8080 $(DOCKER_IMAGE)

swagger:
	swag init -g cmd/server/main.go -o docs

.PHONY: all build test clean run docker-build docker-run swagger
