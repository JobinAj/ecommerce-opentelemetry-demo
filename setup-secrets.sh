#!/bin/bash

# usage: ./setup-secrets.sh <db_password>

if [ -z "$1" ]; then
    echo "Usage: ./setup-secrets.sh <db_password>"
    exit 1
fi

DB_PASSWORD=$1

echo "🔐 Creating 'db-credentials' Secret in 'apps' namespace..."
kubectl create secret generic db-credentials \
  --namespace apps \
  --from-literal=password="$DB_PASSWORD" \
  --dry-run=client -o yaml | kubectl apply -f -

echo "✅ Secret created!"
