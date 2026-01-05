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

// CartItem represents an item in a cart
type CartItem struct {
	ProductID     string  `json:"productId"`
	ProductName   string  `json:"productName"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	SelectedSize  string  `json:"selectedSize"`
	SelectedColor string  `json:"selectedColor"`
}

// Cart represents a shopping cart
type Cart struct {
	ID        string      `json:"id"`
	UserID    string      `json:"userId"`
	Items     []CartItem  `json:"items"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// Order represents an order
type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"userId"`
	Items     []CartItem  `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// CreateCart creates a new cart in the database
func CreateCart(userID string) (*Cart, error) {
	cartID := fmt.Sprintf("cart_%d", time.Now().UnixNano())
	now := time.Now()

	// Begin transaction
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert the cart
	_, err = tx.Exec(`
		INSERT INTO carts (id, user_id, total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`,
		cartID, userID, 0.0, now, now)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Cart{
		ID:        cartID,
		UserID:    userID,
		Items:     []CartItem{},
		Total:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetCart retrieves a cart by its ID
func GetCart(cartID string) (*Cart, error) {
	var cart Cart
	var total sql.NullFloat64
	var createdAt, updatedAt time.Time

	err := DB.QueryRow(`
		SELECT id, user_id, total, created_at, updated_at
		FROM carts 
		WHERE id = $1`, cartID).Scan(
		&cart.ID,
		&cart.UserID,
		&total,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, err
	}

	cart.Total = total.Float64
	cart.CreatedAt = createdAt
	cart.UpdatedAt = updatedAt

	// Get cart items
	cart.Items, err = getCartItems(cartID)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

// getCartItems retrieves all items for a given cart
func getCartItems(cartID string) ([]CartItem, error) {
	rows, err := DB.Query(`
		SELECT product_id, product_name, price, quantity, selected_size, selected_color
		FROM cart_items 
		WHERE cart_id = $1
		ORDER BY created_at`, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		var quantity sql.NullInt64
		var price sql.NullFloat64

		err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&price,
			&quantity,
			&item.SelectedSize,
			&item.SelectedColor,
		)
		if err != nil {
			return nil, err
		}

		item.Price = price.Float64
		item.Quantity = int(quantity.Int64)

		items = append(items, item)
	}

	return items, nil
}

// AddItemToCart adds an item to the cart
func AddItemToCart(cartID string, item CartItem) error {
	// Begin transaction
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert the cart item
	_, err = tx.Exec(`
		INSERT INTO cart_items (cart_id, product_id, product_name, price, quantity, selected_size, selected_color)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		cartID, item.ProductID, item.ProductName, item.Price, item.Quantity, item.SelectedSize, item.SelectedColor)
	if err != nil {
		return err
	}

	// Update the cart total
	total, err := calculateCartTotal(cartID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE carts SET total = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`, total, cartID)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// RemoveItemFromCart removes an item from the cart
func RemoveItemFromCart(cartID, productID string) error {
	// Begin transaction
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete the cart item
	_, err = tx.Exec(`
		DELETE FROM cart_items 
		WHERE cart_id = $1 AND product_id = $2`, cartID, productID)
	if err != nil {
		return err
	}

	// Update the cart total
	total, err := calculateCartTotal(cartID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE carts SET total = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`, total, cartID)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}

// calculateCartTotal calculates the total price of items in a cart
func calculateCartTotal(cartID string) (float64, error) {
	var total float64
	err := DB.QueryRow(`
		SELECT COALESCE(SUM(price * quantity), 0)
		FROM cart_items
		WHERE cart_id = $1`, cartID).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// CreateOrder creates an order from a cart
func CreateOrder(cartID string) (*Order, error) {
	// Get the cart
	cart, err := GetCart(cartID)
	if err != nil {
		return nil, err
	}

	orderID := fmt.Sprintf("ORD_%d", time.Now().UnixNano())
	now := time.Now()

	// Begin transaction
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert the order
	_, err = tx.Exec(`
		INSERT INTO orders (id, user_id, total, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		orderID, cart.UserID, cart.Total, "pending", now, now)
	if err != nil {
		return nil, err
	}

	// Insert order items
	for _, item := range cart.Items {
		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, product_id, product_name, price, quantity, selected_size, selected_color)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			orderID, item.ProductID, item.ProductName, item.Price, item.Quantity, item.SelectedSize, item.SelectedColor)
		if err != nil {
			return nil, err
		}
	}

	// Delete the cart and its items
	_, err = tx.Exec("DELETE FROM cart_items WHERE cart_id = $1", cartID)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec("DELETE FROM carts WHERE id = $1", cartID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Order{
		ID:        orderID,
		UserID:    cart.UserID,
		Items:     cart.Items,
		Total:     cart.Total,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// GetOrder retrieves an order by its ID
func GetOrder(orderID string) (*Order, error) {
	var order Order
	var total sql.NullFloat64
	var status string
	var createdAt, updatedAt time.Time

	err := DB.QueryRow(`
		SELECT id, user_id, total, status, created_at, updated_at
		FROM orders 
		WHERE id = $1`, orderID).Scan(
		&order.ID,
		&order.UserID,
		&total,
		&status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	order.Total = total.Float64
	order.Status = status
	order.CreatedAt = createdAt
	order.UpdatedAt = updatedAt

	// Get order items
	order.Items, err = getOrderItems(orderID)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// getOrderItems retrieves all items for a given order
func getOrderItems(orderID string) ([]CartItem, error) {
	rows, err := DB.Query(`
		SELECT product_id, product_name, price, quantity, selected_size, selected_color
		FROM order_items 
		WHERE order_id = $1
		ORDER BY created_at`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		var quantity sql.NullInt64
		var price sql.NullFloat64

		err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&price,
			&quantity,
			&item.SelectedSize,
			&item.SelectedColor,
		)
		if err != nil {
			return nil, err
		}

		item.Price = price.Float64
		item.Quantity = int(quantity.Int64)

		items = append(items, item)
	}

	return items, nil
}

// UpdateOrderStatus updates the status of an order
func UpdateOrderStatus(orderID, status string) error {
	_, err := DB.Exec(`
		UPDATE orders 
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`, status, orderID)
	return err
}

// GetUserOrders retrieves all orders for a user
func GetUserOrders(userID string) ([]Order, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, total, status, created_at, updated_at
		FROM orders 
		WHERE user_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		var total sql.NullFloat64
		var status string
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&total,
			&status,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		order.Total = total.Float64
		order.Status = status
		order.CreatedAt = createdAt
		order.UpdatedAt = updatedAt

		// Get order items
		order.Items, err = getOrderItems(order.ID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}