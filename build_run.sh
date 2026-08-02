#!/bin/bash

# Configuration
IMAGE_NAME="devops-interview-service"
CONTAINER_NAME="devops-interview-container"
PORT=3000

echo "Building the Docker container image..."
docker build -t $IMAGE_NAME .

# Check if a container with the same name is already running and remove it
if [ "$(docker ps -aq -f name=$CONTAINER_NAME)" ]; then
    echo "Cleaning up old container..."
    docker rm -f $CONTAINER_NAME
fi

echo "Running the container on port $PORT..."
docker run -d --name $CONTAINER_NAME -p $PORT:$PORT -e PORT=$PORT $IMAGE_NAME

echo "Container is up and running!"