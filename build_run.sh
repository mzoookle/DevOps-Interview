#!/bin/bash
set -euo pipefail

# Set default values for environment variables if not set
PORT=3000
IMAGE_NAME="math-server"
CONTAINER_NAME="math-server-container"

if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

echo "Using PORT=${PORT}, IMAGE_NAME=${IMAGE_NAME}, CONTAINER_NAME=${CONTAINER_NAME}"

# Build the Docker image
echo "Building the Docker container image..."
if ! docker build --no-cache -t "${IMAGE_NAME}" .; then
    echo "Docker build failed. Exiting."
    exit 1
fi

# Check if a container with the same name is already running and remove it
if [ "$(docker ps -aq -f name="${CONTAINER_NAME}")" ]; then
    echo "Cleaning up old container..."
    docker rm -f "${CONTAINER_NAME}"
fi

# Run the Docker container
echo "Running the container on port "${PORT}"..."
container_id=$(docker run -d --name "${CONTAINER_NAME}" -p "${PORT}:${PORT}" -e PORT="${PORT}" "${IMAGE_NAME}")
if [ $? -ne 0 ]; then
    echo "Docker run failed. Exiting." >&2
    exit 1
fi
echo "Container is up and running: ${container_id}"