# Variables
BINARY_NAME=test-psql-conn
BUILD_DIR=bin
MAIN_FILE=main.go
EXT=$(if $(filter windows,$(GOOS)),.exe,)

# Targets
.PHONY: all clean build run lint

all: build

build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@go build -mod=readonly -ldflags "-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)$(EXT) .
	@echo "Build complete."

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)$(EXT)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

lint:
	@./scripts/util/check-lint.sh
