#!/bin/bash

#
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$SCRIPT_DIR")"

# Read the configuration..
source $SCRIPT_DIR/podman-data.sh

# 
FULL_IMAGE_NAME="${REGISTRY_HOST}/$REGISTRY_USER/${IMAGE_NAME}:${IMAGE_TAG}"

if [ -z "$FULL_IMAGE_NAME" ]; then
    echo "FULL_IMAGE_NAME is empty"
    exit 1
fi


# Check if the image exists..
if podman images --format "{{.Repository}}:{{.Tag}}" | grep -q "^${FULL_IMAGE_NAME}$"; then
    echo "\$ podman rmi ${FULL_IMAGE_NAME}"
    podman rmi "${FULL_IMAGE_NAME}"
		echo ""
		echo ""
fi


# Build the Docker image
podman build -t $REGISTRY_HOST/$REGISTRY_USER/$IMAGE_NAME:$IMAGE_TAG --file $PARENT_DIR/Dockerfile $PARENT_DIR

IMAGE_ID=$(podman images --format "{{.Id}}" "$FULL_IMAGE_NAME")
if [ -z "$IMAGE_ID" ]; then
    echo "Image not found: $FULL_IMAGE_NAME"
    exit 1
fi

# Log in to the registry..
podman_login_quay.io.sh

# Push the image to the registry
podman push $FULL_IMAGE_NAME

# Print images
echo
podman images
