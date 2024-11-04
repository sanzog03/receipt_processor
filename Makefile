# Go params
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Main package dir
MAIN=./cmd/server/main.go

# Output binary name
BINARY_NAME=main

# Targets

all: test swagger build

swagger:
	swag init -g cmd/server/main.go --parseDependency --parseInternal

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN)

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GORUN) $(MAIN)

.PHONY: swagger build test clean run