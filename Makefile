# Go params
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Main package name
MAIN=receiptProcessor

# Output binary name
BINARY_NAME=main

# Swagger parameters
SWAGGER_OUT=docs
SWAGGER_DIR=cmd,internal

# Targets
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN)

test:
	$(GOTEST) -v ./tests/...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	$(GORUN)  $(MAIN)
