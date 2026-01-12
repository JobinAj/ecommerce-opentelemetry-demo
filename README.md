# Versace E-Commerce - OpenTelemetry Demo

This repository contains a full-stack e-commerce application designed to demonstrate OpenTelemetry integration across Go microservices and a React frontend.

## Architecture

- **Frontend**: React + Vite (Port 5174)
- **Product Service**: Go (Port 8001)
- **Cart & Order Service**: Go (Port 8002)
- **Payment Service**: Go (Port 8003)
- **Database**: AWS RDS PostgreSQL
- **Payment Gateway Mock**: Stripe-mock (Port 12111)

---

## 🚀 Getting Started

### 1. Prerequisites
- **Go**: 1.21+
- **Node.js**: 18+
- **PostgreSQL Client**: `psql` (for database verification)
- **Stripe Mock**: [stripe-mock](https://github.com/stripe/stripe-mock)

---

### 2. Database Connectivity (AWS RDS)

The application connects to an AWS RDS instance. Ensure you have the following environment variables set when running the backend services:

```bash
export DB_HOST=database-1.ce3s0w06y1xp.us-east-1.rds.amazonaws.com
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD='<FQxX1WK8!:PA4-K~hK:X:Y#zwfQ'
export DB_NAME=postgres
```

---

### 3. Running the Services

#### A. Start Stripe Mock
If not already running, start the Stripe mock server:
```bash
docker run --rm -it -p 12111:12111 -p 12112:12112 stripe/stripe-mock
```

#### B. Run Backend Services
You need to start each service in a separate terminal. Ensure the DB environment variables are exported in each terminal.

**Product Service (Port 8001):**
```bash
cd src/product-service
go run main.go
```

**Cart & Order Service (Port 8002):**
```bash
cd src/cart-order-service
go run main.go
```

**Payment Service (Port 8003):**
```bash
cd src/payment-service
go run main.go
```

#### C. Run Frontend (Port 5174)
The frontend uses Vite. Ensure the backend services are running first.

```bash
cd src/frontend
npm install
npm run dev
```

---

## 💳 Testing the Checkout

For testing the payment flow, use the following test card details:

- **Card Number**: `4242 4242 4242 4242`
- **Expiry**: Any future date (e.g., `12/25`)
- **CVV**: Any 3 digits (e.g., `123`)

---

## 🛠 Troubleshooting

- **CORS Errors**: If you encounter CORS issues, ensure that all backend services are successfully responding to `OPTIONS` preflight requests.
- **Port Conflicts**: If port 5173 is in use, Vite will automatically try 5174. Ensure your `.env` or `config.ts` matches the active port.
- **DB Connection**: Ensure your IP is whitelisted in the AWS RDS security group if running from a local machine (not needed if running in the provided environment).
