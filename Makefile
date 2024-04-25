BINARY_NAME=solve
BUILD_DIR=./cmd/solve

all: build

build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) $(BUILD_DIR)

run: build
	@echo "Running..."
	@./$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm $(BINARY_NAME)

.PHONY: all build run clean
