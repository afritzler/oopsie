GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=oopsie
BINARY_LINUX=$(BINARY_NAME)_linux

all: test build
build:
		$(GOBUILD) -o $(BINARY_NAME) -v main.go
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_LINUX)
lint:
		golangci-lint run
run:
		$(GOBUILD) -o $(BINARY_NAME) -v main.go
		./$(BINARY_NAME)
deps:
		# $(GOGET) ...

# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v

docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/afritzler/oopsie golang:latest go build -o "$(BINARY_NAME)" -v
