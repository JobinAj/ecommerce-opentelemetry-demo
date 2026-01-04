#!/bin/bash

# Migration script for Versace E-Commerce Backend

set -e  # Exit on any error

echo "Starting database migration for Versace E-Commerce Backend..."

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "Error: psql is not installed or not in PATH"
    exit 1
fi

# Get database connection details from environment variables or use defaults
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-ecommerce}

echo "Connecting to PostgreSQL at $DB_HOST:$DB_PORT as user $DB_USER"

# Create the database if it doesn't exist
echo "Creating database $DB_NAME if it doesn't exist..."
createdb -h $DB_HOST -p $DB_PORT -U $DB_USER -w $DB_PASSWORD $DB_NAME || true

# Apply the schema
echo "Applying database schema..."
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f /root/ecommerce-opentelemetry-demo/src/backend/db-schema/schema.sql

echo "Database migration completed successfully!"

# Insert initial product data if the products table is empty
echo "Checking if products table needs initial data..."
PRODUCT_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM products;" | xargs)

if [ "$PRODUCT_COUNT" -eq 0 ]; then
    echo "Inserting initial product data..."
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
    INSERT INTO products (id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock) VALUES
    ('1', 'Silk Baroque Shirt', 'Tops', 1200.00, 'https://images.pexels.com/photos/3622613/pexels-photo-3622613.jpeg', 'Luxurious silk shirt with signature Versace Baroque pattern', 4.8, 124, '{XS, S, M, L, XL}', '{Black, Gold, White}', true),
    ('2', 'Gold Medusa Blazer', 'Outerwear', 2500.00, 'https://images.pexels.com/photos/3622614/pexels-photo-3622614.jpeg', 'Statement blazer featuring iconic Medusa head emblem', 4.9, 89, '{XS, S, M, L, XL}', '{Black, Navy, Charcoal}', true),
    ('3', 'Versace Print T-Shirt', 'Tops', 450.00, 'https://images.pexels.com/photos/3622615/pexels-photo-3622615.jpeg', 'Classic cotton t-shirt with Versace logo print', 4.6, 256, '{XS, S, M, L, XL, XXL}', '{White, Black, Red, Navy}', true),
    ('4', 'Tailored Silk Trousers', 'Bottoms', 1800.00, 'https://images.pexels.com/photos/3622616/pexels-photo-3622616.jpeg', 'Elegant silk trousers with perfect drape', 4.7, 142, '{XS, S, M, L, XL}', '{Black, White, Beige}', true),
    ('5', 'Black Leather Jacket', 'Outerwear', 3200.00, 'https://images.pexels.com/photos/3622617/pexels-photo-3622617.jpeg', 'Premium leather jacket with signature detailing', 4.9, 198, '{XS, S, M, L, XL}', '{Black, Brown}', true),
    ('6', 'Gold Chain Dress', 'Dresses', 2800.00, 'https://images.pexels.com/photos/3622618/pexels-photo-3622618.jpeg', 'Stunning dress with gold chain embellishments', 4.8, 167, '{XS, S, M, L}', '{Black, Gold, Silver}', true),
    ('7', 'Premium Denim Jeans', 'Bottoms', 950.00, 'https://images.pexels.com/photos/3622619/pexels-photo-3622619.jpeg', 'High-quality denim with Versace branding', 4.7, 203, '{24, 25, 26, 27, 28, 29, 30, 31, 32}', '{Dark Blue, Light Blue, Black}', true),
    ('8', 'Silk Evening Gown', 'Dresses', 4500.00, 'https://images.pexels.com/photos/3622620/pexels-photo-3622620.jpeg', 'Breathtaking silk gown for special occasions', 5.0, 87, '{XS, S, M, L}', '{Black, Red, White}', true);
    "
    echo "Initial product data inserted successfully!"
else
    echo "Products table already has data, skipping initial data insertion."
fi

echo "Migration script completed!"