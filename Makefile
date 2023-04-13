# Go parameters
# Set go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Linter variables
# Set linter version and binary location
LINTER_VERSION := v1.52.2
LINTER_BINARY := ./bin/golangci-lint

# Binary name
# Set binary name
BINARY_NAME=Similarity

.PHONY: all
# Target to run tests and build the binary
all: lint mocks test build

# Build the project
.PHONY: build
build:
	$(GOBUILD) -o ./bin/$(BINARY_NAME) -v ./cmd/...

# Run linter
.PHONY: lint
lint: $(LINTER_BINARY)
	$(LINTER_BINARY) run ./...

# Generate mocks
.PHONY: mocks
mocks:
	go generate ./...

# Run tests and generate coverage report
.PHONY: test
test: $(LINTER_BINARY)
	# Run tests
	go test -v -race -coverprofile=coverage.out ./...
	# Generate HTML coverage report
	go tool cover -html=coverage.out -o coverage.html

# Remove binary and linter binary
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf bin/

# Build and run the binary
.PHONY: run
run: build
	./bin/$(BINARY_NAME)

# Download dependencies
.PHONY: deps
deps:
	$(GOGET) ./...

# Install linter if necessary
$(LINTER_BINARY):
	@mkdir -p $(@D)
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
        | sh -s -- -b $(dir $@) $(LINTER_VERSION)
