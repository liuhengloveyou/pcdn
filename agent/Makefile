VERSION=0.0.9
BUILD_TIME=$(shell date "+%FZ%T")
COMMIT_SHA1=$(shell git rev-parse HEAD)
GIT_TAG=$(shell git describe --abbrev=0 --always)

flags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.CommitID=${COMMIT_SHA1} -X main.GitTag=${GIT_TAG}"

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=pcdnagent

all: build

build:
	$(GOBUILD) --ldflags ${flags} -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Cross compilation
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) --ldflags ${flags} -o $(BINARY_NAME) -v

arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) --ldflags ${flags} -o $(BINARY_NAME) -v

