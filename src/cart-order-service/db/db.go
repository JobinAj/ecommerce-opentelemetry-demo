package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password,omitempty"`
	Name         string `json:"name"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// CreateUser creates a new user
func CreateUser(user User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name)
		VALUES ($1, $2, $3, $4)`

	_, err := DB.Exec(query, user.ID, user.Email, user.PasswordHash, user.Name)
	return err
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE email = $1`

	var user User
	err := DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}

type CartItem struct {
	ProductID     string  `json:"productId"`
	ProductName   string  `json:"productName"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	SelectedSize  string  `json:"selectedSize"`
	SelectedColor string  `json:"selectedColor"`
}

type Cart struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
}

type Order struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	Status    string     `json:"status"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
}

// CreateCart creates a new cart for a user
func CreateCart(userID string) (*Cart, error) {
	cartID := generateID()
	query := `
		INSERT INTO carts (id, user_id)
		VALUES ($1, $2)
		RETURNING id, user_id, total, created_at, updated_at`

	var cart Cart
	err := DB.QueryRow(query, cartID, userID).Scan(
		&cart.ID, &cart.UserID, &cart.Total, &cart.CreatedAt, &cart.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	cart.Items = []CartItem{}
	return &cart, nil
}

// GetCart retrieves a cart by its ID
func GetCart(cartID string) (*Cart, error) {
	query := `SELECT id, user_id, total, created_at, updated_at FROM carts WHERE id = $1`

	var cart Cart
	err := DB.QueryRow(query, cartID).Scan(
		&cart.ID, &cart.UserID, &cart.Total, &cart.CreatedAt, &cart.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, err
	}

	// Get cart items
	itemsQuery := `
		SELECT product_id, product_name, price, quantity, selected_size, selected_color
		FROM cart_items
		WHERE cart_id = $1
		ORDER BY created_at`

	rows, err := DB.Query(itemsQuery, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		err := rows.Scan(
			&item.ProductID, &item.ProductName, &item.Price,
			&item.Quantity, &item.SelectedSize, &item.SelectedColor,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	cart.Items = items
	return &cart, nil
}

// AddItemToCart adds an item to the cart
func AddItemToCart(cartID string, item CartItem) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, product_name, price, quantity, selected_size, selected_color)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := DB.Exec(query, cartID, item.ProductID, item.ProductName, item.Price, item.Quantity, item.SelectedSize, item.SelectedColor)
	if err != nil {
		return err
	}

	// Update cart total
	return updateCartTotal(cartID)
}

// RemoveItemFromCart removes an item from the cart
func RemoveItemFromCart(cartID, productID string) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2`
	_, err := DB.Exec(query, cartID, productID)
	if err != nil {
		return err
	}

	// Update cart total
	return updateCartTotal(cartID)
}

// CreateOrder creates an order from a cart
func CreateOrder(cartID string) (*Order, error) {
	// First, get the cart
	cart, err := GetCart(cartID)
	if err != nil {
		return nil, err
	}

	orderID := generateID()

	// Begin transaction
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert order
	orderQuery := `
		INSERT INTO orders (id, user_id, total, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, total, status, created_at, updated_at`

	var order Order
	err = tx.QueryRow(orderQuery, orderID, cart.UserID, cart.Total, "pending").Scan(
		&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Insert order items
	for _, item := range cart.Items {
		itemQuery := `
			INSERT INTO order_items (order_id, product_id, product_name, price, quantity, selected_size, selected_color)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err = tx.Exec(itemQuery, order.ID, item.ProductID, item.ProductName, item.Price, item.Quantity, item.SelectedSize, item.SelectedColor)
		if err != nil {
			return nil, err
		}
	}

	// Clear cart items
	_, err = tx.Exec("DELETE FROM cart_items WHERE cart_id = $1", cartID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Get the order with items
	order.Items = cart.Items
	return &order, nil
}

// GetOrder retrieves an order by its ID
func GetOrder(orderID string) (*Order, error) {
	query := `SELECT id, user_id, total, status, created_at, updated_at FROM orders WHERE id = $1`

	var order Order
	err := DB.QueryRow(query, orderID).Scan(
		&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	// Get order items
	itemsQuery := `
		SELECT product_id, product_name, price, quantity, selected_size, selected_color
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at`

	rows, err := DB.Query(itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		err := rows.Scan(
			&item.ProductID, &item.ProductName, &item.Price,
			&item.Quantity, &item.SelectedSize, &item.SelectedColor,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}

// UpdateOrderStatus updates the status of an order
func UpdateOrderStatus(orderID, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := DB.Exec(query, status, orderID)
	return err
}

// GetUserOrders retrieves all orders for a user
func GetUserOrders(userID string) ([]Order, error) {
	query := `SELECT id, user_id, total, status, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.Total, &order.Status, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

// updateCartTotal recalculates and updates the cart total
func updateCartTotal(cartID string) error {
	query := `
		UPDATE carts
		SET total = (
			SELECT COALESCE(SUM(price * quantity), 0)
			FROM cart_items
			WHERE cart_id = $1
		),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := DB.Exec(query, cartID)
	return err
}

// generateID generates a unique ID
func generateID() string {
	// In a real application, you'd want to use a proper UUID generator
	// For now, we'll use a simple timestamp-based ID
	return fmt.Sprintf("id_%d", time.Now().UnixNano())
}