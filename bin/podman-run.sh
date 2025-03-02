#!/bin/bash

REPOSITORY="localhost"
IMAGE_NAME="links"
TAG="latest"

# Run the container
podman run --rm -it -p 8080:8080 $REPOSITORY/$IMAGE_NAME:$TAG
