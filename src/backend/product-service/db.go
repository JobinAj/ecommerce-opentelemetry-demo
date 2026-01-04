package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	// Get database connection details from environment variables
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "5432")
	user := getEnvOrDefault("DB_USER", "postgres")
	password := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbname := getEnvOrDefault("DB_NAME", "ecommerce")

	// Construct the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database!")
}

// getEnvOrDefault returns the environment variable value or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// Product represents a product in the system
type Product struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Price    float64  `json:"price"`
	Image    string   `json:"image"`
	Desc     string   `json:"description"`
	Rating   float64  `json:"rating"`
	Reviews  int      `json:"reviews"`
	Sizes    []string `json:"sizes"`
	Colors   []string `json:"colors"`
	InStock  bool     `json:"inStock"`
}

// GetProducts retrieves all products from the database
func GetProducts(category, search string) ([]Product, error) {
	var products []Product
	var rows *sql.Rows
	var err error

	if category != "" && search != "" {
		rows, err = DB.Query(`
			SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock
			FROM products 
			WHERE category = $1 AND (name ILIKE $2 OR description ILIKE $2)
			ORDER BY name`, category, "%"+search+"%")
	} else if category != "" {
		rows, err = DB.Query(`
			SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock
			FROM products 
			WHERE category = $1
			ORDER BY name`, category)
	} else if search != "" {
		rows, err = DB.Query(`
			SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock
			FROM products 
			WHERE name ILIKE $1 OR description ILIKE $1
			ORDER BY name`, "%"+search+"%")
	} else {
		rows, err = DB.Query(`
			SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock
			FROM products 
			ORDER BY name`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		var sizesBytes, colorsBytes []byte

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Price,
			&product.Image,
			&product.Desc,
			&product.Rating,
			&product.Reviews,
			&sizesBytes,
			&colorsBytes,
			&product.InStock,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal the array fields
		if sizesBytes != nil {
			if err := json.Unmarshal(sizesBytes, &product.Sizes); err != nil {
				return nil, err
			}
		}
		if colorsBytes != nil {
			if err := json.Unmarshal(colorsBytes, &product.Colors); err != nil {
				return nil, err
			}
		}

		products = append(products, product)
	}

	return products, nil
}

// GetProductByID retrieves a product by its ID
func GetProductByID(id string) (*Product, error) {
	var product Product
	var sizesBytes, colorsBytes []byte

	err := DB.QueryRow(`
		SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock
		FROM products 
		WHERE id = $1`, id).Scan(
		&product.ID,
		&product.Name,
		&product.Category,
		&product.Price,
		&product.Image,
		&product.Desc,
		&product.Rating,
		&product.Reviews,
		&sizesBytes,
		&colorsBytes,
		&product.InStock,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	// Unmarshal the array fields
	if sizesBytes != nil {
		if err := json.Unmarshal(sizesBytes, &product.Sizes); err != nil {
			return nil, err
		}
	}
	if colorsBytes != nil {
		if err := json.Unmarshal(colorsBytes, &product.Colors); err != nil {
			return nil, err
		}
	}

	return &product, nil
}

// GetCategories retrieves all unique categories
func GetCategories() ([]string, error) {
	rows, err := DB.Query("SELECT DISTINCT category FROM products ORDER BY category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// UpdateStock updates the stock status of a product
func UpdateStock(productID string, quantity int) error {
	// For this implementation, we'll just update the in_stock status
	// In a real system, you might want to track actual inventory quantities
	_, err := DB.Exec("UPDATE products SET in_stock = $2 WHERE id = $1", productID, quantity > 0)
	return err
}

// SearchProducts searches for products based on query and price range
func SearchProducts(query, minPrice, maxPrice string) ([]Product, error) {
	var products []Product
	var rows *sql.Rows
	var err error

	// Set default min and max prices if not provided
	min := 0.0
	max := 100000.0

	if minPrice != "" {
		fmt.Sscanf(minPrice, "%f", &min)
	}
	if maxPrice != "" {
		fmt.Sscanf(maxPrice, "%f", &max)
	}

	if query != "" {
		rows, err = DB.Query(`
			SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock
			FROM products 
			WHERE (name ILIKE $1 OR description ILIKE $1) 
			AND price >= $2 AND price <= $3
			ORDER BY name`, "%"+query+"%", min, max)
	} else {
		rows, err = DB.Query(`
			SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock
			FROM products 
			WHERE price >= $1 AND price <= $2
			ORDER BY name`, min, max)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		var sizesBytes, colorsBytes []byte

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Category,
			&product.Price,
			&product.Image,
			&product.Desc,
			&product.Rating,
			&product.Reviews,
			&sizesBytes,
			&colorsBytes,
			&product.InStock,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal the array fields
		if sizesBytes != nil {
			if err := json.Unmarshal(sizesBytes, &product.Sizes); err != nil {
				return nil, err
			}
		}
		if colorsBytes != nil {
			if err := json.Unmarshal(colorsBytes, &product.Colors); err != nil {
				return nil, err
			}
		}

		products = append(products, product)
	}

	return products, nil
}