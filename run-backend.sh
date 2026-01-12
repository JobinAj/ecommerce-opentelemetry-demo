#!/bin/bash

# Load environment variables from .env file
set -a
source /root/ecommerce-opentelemetry-demo/.env
set +a

echo "Environment variables loaded:"
echo "DB_HOST: $DB_HOST"
echo "DB_PORT: $DB_PORT"
echo "DB_USER: $DB_USER"
echo "DB_NAME: $DB_NAME"

# Navigate to the src directory
cd /root/ecommerce-opentelemetry-demo/src

# Build the services
echo "Building backend services..."
go build -o cart-service ./cart-order-service/main.go
go build -o product-service ./product-service/main.go
go build -o payment-service ./payment-service/main.go

echo "Starting backend services..."

# Start the services in the background
echo "Starting Product Service on port 8001..."
./product-service &
PRODUCT_PID=$!

sleep 2

echo "Starting Cart & Order Service on port 8002..."
./cart-service &
CART_PID=$!

sleep 2

echo "Starting Payment Service on port 8003..."
./payment-service &
PAYMENT_PID=$!

sleep 2

# Check if services are running
echo "Checking if services are running..."
if ps -p $PRODUCT_PID > /dev/null; then
    echo "Product Service is running (PID: $PRODUCT_PID)"
else
    echo "Product Service failed to start"
fi

if ps -p $CART_PID > /dev/null; then
    echo "Cart & Order Service is running (PID: $CART_PID)"
else
    echo "Cart & Order Service failed to start"
fi

if ps -p $PAYMENT_PID > /dev/null; then
    echo "Payment Service is running (PID: $PAYMENT_PID)"
else
    echo "Payment Service failed to start"
fi

# Keep the script running
wait $PRODUCT_PID $CART_PID $PAYMENT_PID