# Variables
BINARY_NAME=test-psql-conn
BUILD_DIR=bin
MAIN_FILE=main.go

# Targets
.PHONY: all clean build run

all: build

build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -ldflags "-w -s" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete."

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."