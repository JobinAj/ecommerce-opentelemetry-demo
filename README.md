# Versace E-Commerce - OpenTelemetry Demo

A full-stack e-commerce application demonstrating **OpenTelemetry observability**, **feature flags**, and **Kubernetes deployment** on AWS EKS.

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              AWS EKS Cluster                            │
├─────────────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐    │
│  │  Frontend   │  │  Product    │  │  Cart/Order │  │  Payment    │    │
│  │  (React)    │  │  Service    │  │  Service    │  │  Service    │    │
│  │  Port 3000  │  │  Port 8001  │  │  Port 8002  │  │  Port 8003  │    │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘    │
│         │                │                │                │            │
│         └────────────────┴────────────────┴────────────────┘            │
│                                    │                                    │
│  ┌─────────────────────────────────┼─────────────────────────────────┐  │
│  │                     Observability Stack                           │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐               │  │
│  │  │    OTel     │  │   Jaeger    │  │ Prometheus  │               │  │
│  │  │  Collector  │──│   (Traces)  │  │  (Metrics)  │               │  │
│  │  └─────────────┘  └─────────────┘  └──────┬──────┘               │  │
│  │                                           │                       │  │
│  │                                    ┌──────┴──────┐               │  │
│  │                                    │   Grafana   │               │  │
│  │                                    │ (Dashboards)│               │  │
│  │                                    └─────────────┘               │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                         │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                     │
│  │    flagd    │  │ Stripe-Mock │  │   Locust    │                     │
│  │ (Feature    │  │ (Payment    │  │   (Load     │                     │
│  │   Flags)    │  │  Gateway)   │  │  Testing)   │                     │
│  └─────────────┘  └─────────────┘  └─────────────┘                     │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                           ┌────────┴────────┐
                           │   AWS RDS       │
                           │  PostgreSQL     │
                           └─────────────────┘
```

## 📦 Components

| Component | Technology | Port | Description |
|-----------|------------|------|-------------|
| Frontend | React + Vite | 3000 | Luxury e-commerce UI |
| Product Service | Go | 8001 | Product catalog API |
| Cart/Order Service | Go | 8002 | Shopping cart & order management |
| Payment Service | Go | 8003 | Payment processing (Stripe) |
| Database | PostgreSQL (RDS) | 5432 | Persistent data storage |
| Stripe Mock | Docker | 12111 | Payment gateway simulator |
| flagd | OpenFeature | 8013 | Feature flag service |
| OTel Collector | OpenTelemetry | 4317/4318 | Telemetry aggregation |
| Jaeger | CNCF | 16686 | Distributed tracing UI |
| Prometheus | CNCF | 9090 | Metrics collection |
| Grafana | Grafana | 3000 | Observability dashboards |
| Locust | Python | 8089 | Load testing |

---

## 🚀 Kubernetes Deployment

### Prerequisites

- AWS EKS Cluster
- `kubectl` configured for your cluster
- Docker Hub account (for image pushing)
- AWS RDS PostgreSQL instance

### Step 1: Deploy Infrastructure (Terraform)

```bash
cd terraform
terraform init
terraform apply
```

### Step 2: Install OpenTelemetry Operator

```bash
# Install cert-manager (required for OTel Operator)
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.4/cert-manager.yaml
kubectl wait --for=condition=Available deployment --all -n cert-manager --timeout=300s

# Install OTel Operator
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/latest/download/opentelemetry-operator.yaml
kubectl wait --for=condition=Available deployment --all -n opentelemetry-operator-system --timeout=300s
```

### Step 3: Create Namespaces

```bash
kubectl create namespace apps
kubectl create namespace observability
```

### Step 4: Deploy Database Credentials

```bash
kubectl create secret generic db-credentials \
  --from-literal=password='YOUR_DB_PASSWORD' \
  -n apps
```

### Step 5: Deploy Observability Stack

```bash
kubectl apply -f kubernetes/observability/observability.yaml
```

### Step 6: Deploy Applications

```bash
# Deploy feature flags
kubectl apply -f kubernetes/apps/flagd.yaml

# Deploy addons (stripe-mock)
kubectl apply -f kubernetes/addons/addons.yaml

# Deploy microservices
kubectl apply -f kubernetes/apps/microservices.yaml

# Deploy frontend
kubectl apply -f kubernetes/apps/frontend.yaml

# Deploy load generator
kubectl apply -f kubernetes/apps/locust.yaml
```

### Step 7: Initialize Database

```bash
kubectl apply -f kubernetes/apps/db-init-job.yaml
```

### Step 8: Verify Deployment

```bash
# Check all pods are running
kubectl get pods -n apps
kubectl get pods -n observability

# Check services
kubectl get svc -n apps
kubectl get svc -n observability
```

---

## 📊 Accessing Observability Tools

### Port Forward to Access UIs

```bash
# Jaeger (Tracing)
kubectl port-forward svc/jaeger -n observability 16686:16686

# Grafana (Dashboards)
kubectl port-forward svc/grafana -n observability 3000:3000

# Prometheus (Metrics)
kubectl port-forward svc/prometheus -n observability 9090:9090

# Locust (Load Testing)
kubectl port-forward svc/locust -n apps 8089:8089
```

### Default Credentials

- **Grafana**: admin / admin

---

## 🎛️ Feature Flags

The application uses **OpenFeature** with **flagd** for feature flag management.

### Available Flags

| Flag Name | Description | Default |
|-----------|-------------|---------|
| `paymentServiceFailure` | Simulates payment service failure | off |
| `cartServiceFailure` | Simulates cart service failure | off |
| `productCatalogFailure` | Simulates product catalog failure | off |

### Enabling a Feature Flag

1. Edit `kubernetes/apps/flagd.yaml`:
```yaml
"paymentServiceFailure": {
  "state": "ENABLED",
  "variants": {
    "on": true,
    "off": false
  },
  "defaultVariant": "on"  # Change to "on" to enable
},
```

2. Apply and restart:
```bash
kubectl apply -f kubernetes/apps/flagd.yaml
kubectl rollout restart deployment/flagd -n apps
```

3. Verify the flag is working by attempting a checkout - you should see "Simulated Payment Service Failure".

---

## 🔧 Building & Pushing Docker Images

### Rebuild All Services

```bash
./rebuild_services.sh
```

### Manual Build

```bash
# Build individual service
cd src/cart-order-service
docker build -t yourusername/cart-order-service:latest .
docker push yourusername/cart-order-service:latest

# Restart deployment to pull new image
kubectl rollout restart deployment/cart-order-service -n apps
```

---

## 💳 Testing the Application

### Test Card Details

| Field | Value |
|-------|-------|
| Card Number | 4242 4242 4242 4242 |
| Expiry | Any future date (e.g., 12/25) |
| CVV | Any 3 digits (e.g., 123) |

### Load Testing with Locust

1. Port forward: `kubectl port-forward svc/locust -n apps 8089:8089`
2. Open http://localhost:8089
3. Configure users and spawn rate
4. Start the test

---

## 🔍 Troubleshooting

### Check Pod Logs

```bash
kubectl logs -n apps -l app=payment-service --tail=50
kubectl logs -n observability -l app=otel-collector --tail=50
```

### Verify Traces in Jaeger

```bash
# Query Jaeger API for services
kubectl run curl-check --image=curlimages/curl --restart=Never --rm -i -- \
  curl -s http://jaeger.observability:16686/api/services
```

### Common Issues

| Issue | Solution |
|-------|----------|
| Pods in CrashLoopBackOff | Check logs: `kubectl logs -n apps <pod-name>` |
| No traces in Jaeger | Verify OTel Collector is running and OTLP endpoint is correct |
| Feature flags not working | Ensure flagd service is named `otel-flagd` (not `flagd`) to avoid env var collision |

---

## 📁 Project Structure

```
ecommerce-opentelemetry-demo/
├── kubernetes/
│   ├── apps/
│   │   ├── microservices.yaml    # Go backend services
│   │   ├── frontend.yaml         # React frontend
│   │   ├── flagd.yaml            # Feature flag service
│   │   ├── locust.yaml           # Load testing
│   │   └── db-init-job.yaml      # Database initialization
│   ├── addons/
│   │   └── addons.yaml           # Stripe mock
│   └── observability/
│       └── observability.yaml    # Jaeger, Prometheus, Grafana, OTel Collector
├── src/
│   ├── frontend/                 # React application
│   ├── product-service/          # Go product API
│   ├── cart-order-service/       # Go cart/order API
│   ├── payment-service/          # Go payment API
│   ├── database/                 # SQL schema & seeds
│   └── flagd/                    # Feature flag config
├── terraform/                    # AWS infrastructure
└── rebuild_services.sh           # Docker build script
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

---

## 📄 License

MIT License - See LICENSE file for details.
