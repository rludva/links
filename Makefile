# Name of the application
APP_NAME = links
OUTPUT_DIR = ./build

# Default target
all: build

# Build the application
build:
	if [ ! -d "$(OUTPUT_DIR)" ]; then \
		mkdir -p "$(OUTPUT_DIR)"; \
	fi
	go build -o $(OUTPUT_DIR)/$(APP_NAME) links.go

# Clean up build artifacts
clean:
	rm -rf $(OUTPUT_DIR)

# Run the application
run: build
	$(OUTPUT_DIR)/$(APP_NAME)

.PHONY: all build clean run 
