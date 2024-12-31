# Name of the application
APP_NAME ?= links
OUTPUT_DIR = ./build

# Default output file name..
BUILD_OUTPUT = ${OUTPUT_DIR}/${APP_NAME}

# Image varaibles
IMAGE_VERSION ?= 1.0.1
IMAGE_REPOSITORY ?= quay.io/rludva/$(APP_NAME)

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
	rm -rf $(BUILD_OUTPUT)

# Run the application
run: build
	$(OUTPUT_DIR)/$(APP_NAME)

# Build the application image from the Dockerfile..
docker-build:
	podman build -t $(IMAGE_REPOSITORY)/$(APP_NAME):$(IMAGE_VERSION) .
	podman images

docker-run:
	podman run --rm -it $(IMAGE_REPOSITORY)/$(APP_NAME):$(IMAGE_VERSION) bash

# Open the application in the browser
open:
	xdg-open http://localhost:8080

.PHONY: all build clean run 
