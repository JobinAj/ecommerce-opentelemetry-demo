# Versace E-Commerce - Golang Backend Integration Guide

## Overview

This project includes a production-ready React frontend that can be integrated with the Golang microservices backend located in the `golang-backend/` directory.

## Frontend Architecture

The React frontend is built with:
- **React 18** with TypeScript
- **Tailwind CSS** for styling
- **Lucide React** for icons
- **Vite** for build tooling

### Features Implemented

1. **Product Catalog**
   - Browse all products with categories
   - Search functionality
   - Product filtering by category
   - Detailed product view with ratings and reviews

2. **Shopping Cart**
   - Add/remove items
   - Adjust quantities
   - Real-time total calculation
   - Persistent cart state

3. **Checkout Flow**
   - Secure payment form
   - Address entry
   - Order summary
   - Payment processing

## Golang Microservices

Three microservices are provided in `golang-backend/`:

### 1. Product Service (Port 8001)
**Endpoints:**
- `GET /api/products` - Get all products with optional filtering
- `GET /api/products/{id}` - Get product details
- `GET /api/categories` - Get all categories
- `GET /api/search?q=query&minPrice=0&maxPrice=5000` - Search products
- `POST /api/stock/update` - Update product stock

**Query Parameters for /api/products:**
```
?category=Tops&search=Silk
```

### 2. Cart & Order Service (Port 8002)
**Endpoints:**
- `POST /api/carts` - Create a new cart
- `GET /api/carts/{cartId}` - Get cart details
- `POST /api/carts/{cartId}/items` - Add item to cart
- `DELETE /api/carts/{cartId}/items/{productId}` - Remove item from cart
- `POST /api/orders` - Create order from cart (body: `{cartId: "..."}`)
- `GET /api/orders/{orderId}` - Get order details
- `PUT /api/orders/{orderId}/status` - Update order status
- `GET /api/users/{userId}/orders` - Get user's orders

### 3. Payment Service (Port 8003)
**Endpoints:**
- `POST /api/payments` - Process payment
- `GET /api/payments/{paymentId}` - Get payment details
- `GET /api/payments/order/{orderId}` - Get payment for order
- `POST /api/payments/{paymentId}/refund` - Refund payment

**Payment Request Body:**
```json
{
  "orderId": "ORD_123",
  "amount": 1200.00,
  "cardNumber": "4532123456789010",
  "cardHolder": "John Doe",
  "expiryDate": "12/25",
  "cvv": "123"
}
```

## Integration Steps

### 1. Setup Golang Services

```bash
cd golang-backend

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

### 2. Update Frontend API Calls

To integrate the frontend with your Golang backend, update the API calls in the React components:

**Example - Create a client service:**

```typescript
// src/services/api.ts
const API_BASE = {
  products: 'http://localhost:8001/api',
  cart: 'http://localhost:8002/api',
  payment: 'http://localhost:8003/api',
};

export const getProducts = async () => {
  const response = await fetch(`${API_BASE.products}/products`);
  return response.json();
};

export const processPayment = async (paymentData) => {
  const response = await fetch(`${API_BASE.payment}/payments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(paymentData),
  });
  return response.json();
};
```

### 3. CORS Configuration

All microservices have CORS enabled to accept requests from the frontend. No additional configuration needed during development.

## Testing Endpoints

### Health Check
```bash
curl http://localhost:8001/health
curl http://localhost:8002/health
curl http://localhost:8003/health
```

### Get All Products
```bash
curl http://localhost:8001/api/products
```

### Search Products
```bash
curl "http://localhost:8001/api/search?q=Silk&minPrice=1000&maxPrice=2000"
```

### Process Payment
```bash
curl -X POST http://localhost:8003/api/payments \
  -H "Content-Type: application/json" \
  -d '{
    "orderId": "ORD_1",
    "amount": 1200,
    "cardNumber": "4532123456789010",
    "cardHolder": "John Doe",
    "expiryDate": "12/25",
    "cvv": "123"
  }'
```

## Deployment

### Frontend
The React frontend is built with:
```bash
npm run build
```

Deploy the `dist/` folder to:
- Vercel
- Netlify
- GitHub Pages
- AWS S3 + CloudFront
- Any static hosting service

### Golang Services

Deploy the Golang services to:
- **Docker**: Containerize each service
- **Kubernetes**: Deploy as separate pods
- **AWS EC2/ECS**: Run as services
- **Heroku/Railway**: Deploy directly
- **Digital Ocean/Linode**: VPS deployment

**Example Dockerfile:**
```dockerfile
FROM golang:1.21-alpine as builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o service main.go

FROM alpine:latest
COPY --from=builder /app/service /service
EXPOSE 8001
CMD ["/service"]
```

## Service Communication

Example of how Cart Service can call Product Service to verify stock:

```go
// In cart-order-service
resp, _ := http.Get("http://localhost:8001/api/products/" + productID)
defer resp.Body.Close()
json.NewDecoder(resp.Body).Decode(&product)
```

## Data Persistence

Currently, all microservices store data in memory. To add persistence:

### Add Database Layer

```go
import "database/sql"
import _ "github.com/lib/pq"

db, err := sql.Open("postgres", "postgres://user:password@localhost/versace")
```

### Example: Get products from database
```go
func getAllProducts(db *sql.DB) ([]Product, error) {
    rows, err := db.Query("SELECT id, name, price, category FROM products")
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var p Product
        rows.Scan(&p.ID, &p.Name, &p.Price, &p.Category)
        products = append(products, p)
    }
    return products, nil
}
```

## Running in Development

### Option 1: Local Development
```bash
# Frontend (one terminal)
npm run dev

# Backend Services (separate terminals)
cd golang-backend/product-service && go run main.go
cd golang-backend/cart-order-service && go run main.go
cd golang-backend/payment-service && go run main.go
```

### Option 2: Docker Compose
Create `docker-compose.yml`:
```yaml
version: '3'
services:
  frontend:
    build: .
    ports:
      - "5173:5173"

  product-service:
    build: ./golang-backend/product-service
    ports:
      - "8001:8001"

  cart-service:
    build: ./golang-backend/cart-order-service
    ports:
      - "8002:8002"

  payment-service:
    build: ./golang-backend/payment-service
    ports:
      - "8003:8003"
```

Run with: `docker-compose up`

## Security Considerations

1. **Payment Data**: Never store full card details. The frontend sends to backend; backend should never log or store them.
2. **CORS**: In production, restrict CORS to your domain only
3. **Authentication**: Add JWT token validation for secure endpoints
4. **HTTPS**: Always use HTTPS in production
5. **Rate Limiting**: Implement rate limiting on payment endpoints
6. **Input Validation**: Validate all inputs (card number, amounts, etc.)

## Next Steps

1. Add user authentication (JWT)
2. Integrate with real payment provider (Stripe, PayPal)
3. Add database persistence (PostgreSQL)
4. Implement logging and monitoring
5. Add unit and integration tests
6. Set up CI/CD pipeline
7. Deploy to production infrastructure
