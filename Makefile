.DEFAULT_GOAL := move

BIN_DIR=./bin
BINARY_NAME=rummi-q
CONFIG_FILE=config.yaml

# Default OS and Architecture compiler target
GOOS=darwin
GOARCH=arm64

.PHONY: build clean move run time vet dep

time:
	@date

clean: time
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)/*
	@mkdir -p $(BIN_DIR)

# Get dependencies
dep:
	@echo "Checking dependencies..."
	@go mod tidy

vet:
	@echo "Running go vet..."
	@# Vet all packages except tmp
	@go vet $(shell go list ./... | grep -v '/tmp')

build: time clean dep vet
	@echo "Building $(BINARY_NAME)..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ${BINARY_NAME} ./cmd/

move: time build
	@echo "Moving binary to $(BIN_DIR)..."
	@mv $(BINARY_NAME) $(BIN_DIR)/
	@echo "Moving config to $(BIN_DIR)..."
	@cp $(CONFIG_FILE) $(BIN_DIR)/

# Build for MacOS
build-macos: time vet
	@echo "Building $(BINARY_NAME) for MacOS/arm64..."
	@GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME) ./cmd/

# Build for Linux
build-linux: time vet
	@echo "Building $(BINARY_NAME) for linux/amd64..."
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) ./cmd/

# Build for Windows
build-windows: time vet
	@echo "Building $(BINARY_NAME).exe for windows/amd64..."
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe ./cmd/

run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)