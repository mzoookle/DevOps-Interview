#!/bin/bash
set -euo pipefail

if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Check if required environment variables are set
if [ -z "${PORT:-}" ] || [ -z "${NAME:-}" ]; then
    echo "Error: PORT and NAME environment variables must be set in .env file."
    exit 1
fi

echo "Using PORT=${PORT}, NAME=${NAME}"

# Build the Docker image
echo "Building the Docker container image..."
if ! docker build --no-cache -t "${NAME}" .; then
    echo "Docker build failed. Exiting."
    exit 1
fi

# Check if a container with the same name is already running and remove it
if [ "$(docker ps -aq -f name="${NAME}")" ]; then
    echo "Cleaning up old container..."
    docker rm -f "${NAME}"
fi

# Run the Docker container
echo "Running the container on port ${PORT}..."
container_id="$(docker run -d --name "${NAME}" -p "${PORT}:${PORT}" --env-file .env "${NAME}")"
if [ $? -ne 0 ]; then
    echo "Docker run failed. Exiting." >&2
    exit 1
fi
echo "Container is up and running: ${container_id}"