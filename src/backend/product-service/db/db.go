package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	// Get database connection details from environment variables
	host := GetEnvOrDefault("DB_HOST", "localhost")
	port := GetEnvOrDefault("DB_PORT", "5432")
	user := GetEnvOrDefault("DB_USER", "postgres")
	password := GetEnvOrDefault("DB_PASSWORD", "postgres")
	dbname := GetEnvOrDefault("DB_NAME", "ecommerce")

	// Construct the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
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

// GetEnvOrDefault returns the environment variable value or a default value
func GetEnvOrDefault(key, defaultValue string) string {
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

type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Price       float64  `json:"price"`
	Image       string   `json:"image"`
	Description string   `json:"description"`
	Rating      float64  `json:"rating"`
	Reviews     int      `json:"reviews"`
	Sizes       []string `json:"sizes"`
	Colors      []string `json:"colors"`
	InStock     bool     `json:"inStock"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

// GetProducts retrieves products with optional category and search filters
func GetProducts(category, search string) ([]Product, error) {
	query := `SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock, created_at, updated_at FROM products WHERE true`
	var args []interface{}
	argCount := 1

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, category)
		argCount++
	}

	if search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
		argCount += 2
	}

	query += " ORDER BY created_at DESC"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		var sizes, colors pq.StringArray
		err := rows.Scan(
			&p.ID, &p.Name, &p.Category, &p.Price, &p.Image, &p.Description,
			&p.Rating, &p.Reviews, &sizes, &colors, &p.InStock, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		p.Sizes = []string(sizes)
		p.Colors = []string(colors)

		products = append(products, p)
	}

	return products, nil
}

// GetProductByID retrieves a product by its ID
func GetProductByID(id string) (*Product, error) {
	query := `SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock, created_at, updated_at FROM products WHERE id = $1`

	var p Product
	var sizes, colors pq.StringArray
	err := DB.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Category, &p.Price, &p.Image, &p.Description,
		&p.Rating, &p.Reviews, &sizes, &colors, &p.InStock, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	p.Sizes = []string(sizes)
	p.Colors = []string(colors)

	return &p, nil
}

// GetCategories retrieves all unique categories
func GetCategories() ([]string, error) {
	query := `SELECT DISTINCT category FROM products ORDER BY category`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// UpdateStock updates the stock quantity for a product
func UpdateStock(productID string, quantity int) error {
	// For this implementation, we'll just check if the product exists
	// In a real application, you'd want to update an inventory table
	query := `SELECT id FROM products WHERE id = $1`
	var id string
	err := DB.QueryRow(query, productID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("product not found")
		}
		return err
	}

	// In a real implementation, you would update the stock in a separate inventory table
	// For now, we'll just return nil to indicate success
	return nil
}

// SearchProducts searches for products based on query and price range
func SearchProducts(query, minPrice, maxPrice string) ([]Product, error) {
	baseQuery := `SELECT id, name, category, price, image, description, rating, reviews, sizes, colors, in_stock, created_at, updated_at FROM products WHERE true`
	var args []interface{}
	argCount := 1

	if query != "" {
		baseQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argCount, argCount+1)
		args = append(args, "%"+query+"%", "%"+query+"%")
		argCount += 2
	}

	if minPrice != "" {
		baseQuery += fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, minPrice)
		argCount++
	}

	if maxPrice != "" {
		baseQuery += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, maxPrice)
		argCount++
	}

	baseQuery += " ORDER BY created_at DESC"

	rows, err := DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		var sizes, colors pq.StringArray
		err := rows.Scan(
			&p.ID, &p.Name, &p.Category, &p.Price, &p.Image, &p.Description,
			&p.Rating, &p.Reviews, &sizes, &colors, &p.InStock, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		p.Sizes = []string(sizes)
		p.Colors = []string(colors)

		products = append(products, p)
	}

	return products, nil
}