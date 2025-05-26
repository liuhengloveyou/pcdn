VERSION=0.0.1
BUILD_TIME=$(shell date "+%FZ%T")
COMMIT_SHA1=$(shell git rev-parse HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --always)

flags="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1} -X main.GitTag=${GIT_TAG}"

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=agent
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_ARM=$(BINARY_NAME)_arm

all: build

build:
	$(GOBUILD) --ldflags ${flags} -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_ARM)

# Cross compilation
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) --ldflags ${flags} -o $(BINARY_UNIX) -v

arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) --ldflags ${flags} -o $(BINARY_ARM) -v
