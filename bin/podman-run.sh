#!/bin/bash

#
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PARENT_DIR="$(dirname "$SCRIPT_DIR")"


# Read the configuration..
source $SCRIPT_DIR/podman-data.sh

# Run the container
podman run --rm -it -p 8080:8080 $REGISTRY_HOST/$REGISTRY_USER/$IMAGE_NAME:$IMAGE_TAG
