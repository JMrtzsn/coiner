 # Go parameters
    GOCMD=go
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    BINARY_NAME=coiner
    BINARY_UNIX=$(BINARY_NAME)_unix

    .PHONY: build
    build:
		$(GOBUILD) -o $(BINARY_NAME) -v
    .PHONY: test
    test:
		$(GOTEST) -v ./...
    .PHONY: clean
    clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
    .PHONY: run
    run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

    # Cross compilation
    .PHONY: build-linux
    build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v