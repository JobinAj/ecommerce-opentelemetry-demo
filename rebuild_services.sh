#!/bin/bash
set -e

echo "Building cart-order-service..."
docker build -t jobinaj/cart-order-service:latest src/cart-order-service
docker push jobinaj/cart-order-service:latest

echo "Building product-service..."
docker build -t jobinaj/product-service:latest src/product-service
docker push jobinaj/product-service:latest

echo "Building payment-service..."
docker build -t jobinaj/payment-service:latest src/payment-service
docker push jobinaj/payment-service:latest

echo "All services built and pushed."
