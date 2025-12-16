.DEFAULT_GOAL := move

BIN_DIR=./bin
BINARY_NAME=rummi-q
CONFIG_FILE=config.yaml

# Default OS and Architecture compiler target
GOOS=darwin
GOARCH=arm64

.PHONY: build clean move run time

time:
	@date

clean: time
	@echo "Cleaning..."
	rm -rf $(BIN_DIR)/*
	mkdir -p $(BIN_DIR)

build: time clean
	@echo "Building $(BINARY_NAME)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ${BINARY_NAME} ./cmd/

move: time build
	@echo "Moving binary to $(BIN_DIR)..."
	mv $(BINARY_NAME) $(BIN_DIR)/
	@echo "Moving config to $(BIN_DIR)..."
	cp $(CONFIG_FILE) $(BIN_DIR)/

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)