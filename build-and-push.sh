#!/bin/bash
set -e

# Docker Credentials
if [ -z "$1" ]; then
    echo "Usage: ./build-and-push.sh <docker_pat>"
    exit 1
fi

DOCKER_USER="jobinaj"
DOCKER_PAT=$1

echo "🐳 Logging in to Docker Hub..."
echo "$DOCKER_PAT" | docker login --username "$DOCKER_USER" --password-stdin

# Services to build
SERVICES=("product-service" "payment-service" "cart-order-service")

for service in "${SERVICES[@]}"; do
  echo "🛠️  Building $service..."
  cd src/$service
  docker build -t $DOCKER_USER/$service:latest .
  
  echo "🚀 Pushing $service..."
  docker push $DOCKER_USER/$service:latest
  
  cd ../.. # Return to root
  echo "✅ $service done!"
done

echo "🎉 All images built and pushed successfully!"
