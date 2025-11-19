# Versace E-Commerce Backend - Golang Microservices

This directory contains three microservices for the Versace e-commerce platform:

## Microservices Architecture

### 1. Product Service
- **Port**: 8001
- **Responsibilities**:
  - Product catalog management
  - Category management
  - Product filtering and search
  - Inventory tracking

### 2. Cart & Order Service
- **Port**: 8002
- **Responsibilities**:
  - Shopping cart management
  - Order creation and management
  - Order history retrieval
  - Order status updates

### 3. Payment Service
- **Port**: 8003
- **Responsibilities**:
  - Payment processing
  - Payment validation
  - Transaction history
  - Payment status tracking

## Setup Instructions

### Prerequisites
- Go 1.21 or higher
- PostgreSQL (optional, for data persistence)

### Running the Services

```bash
# Terminal 1 - Product Service
cd product-service
go mod download
go run main.go

# Terminal 2 - Cart & Order Service
cd cart-order-service
go mod download
go run main.go

# Terminal 3 - Payment Service
cd payment-service
go mod download
go run main.go
```

## API Integration

Update your React frontend `.env` file to point to these services:

```
VITE_PRODUCT_API=http://localhost:8001/api
VITE_CART_API=http://localhost:8002/api
VITE_PAYMENT_API=http://localhost:8003/api
```

## Communication Between Services

Services communicate via HTTP REST APIs. The Payment Service may need to verify cart contents from the Cart Service.
