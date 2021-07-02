 # Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    GOVET=$(GOCMD) vet
    GOFMT=$(GOCMD) fmt
    GOFIX=$(GOCMD) fix
    GOLINT=golangci-lint run
    BINARY_NAME=coiner
    BINARY_UNIX=$(BINARY_NAME)_unix

    .PHONY: build
    build:
		$(GOBUILD) -o $(BINARY_NAME) -v

    .PHONY: test
    test:
		$(GOTEST) -v ./...

    .PHONY: lint
    lint:
		$(GOLINT) ./...

    .PHONY: run
    run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

	.PHONY: fix
    fix:
		$(GOFIX) ./...

    .PHONY: fmt
    fmt:
		$(GOFMT) ./...

	.PHONY: hygiene
    	hygiene: build fmt fix lint vet test

    .PHONY: vet
    vet:
		$(GOVET) ./...

    # Cross compilation
    .PHONY: build-linux
    build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

     .PHONY: dockerfile
     dockerfile:
		docker build -t $(BINARY_NAME) .
