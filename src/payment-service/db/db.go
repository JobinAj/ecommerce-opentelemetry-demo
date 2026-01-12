package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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

type PaymentRequest struct {
	OrderID    string  `json:"orderId"`
	Amount     float64 `json:"amount"`
	Currency   string  `json:"currency"`
	CardNumber string  `json:"cardNumber"`
	CardHolder string  `json:"cardHolder"`
	ExpiryDate string  `json:"expiryDate"`
	CVV        string  `json:"cvv"`
}

type Payment struct {
	ID            string  `json:"id"`
	OrderID       string  `json:"orderId"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	CardLastFour  string  `json:"cardLastFour"`
	TransactionID string  `json:"transactionId"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

type PaymentResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Payment Payment `json:"payment,omitempty"`
}

// CreatePayment creates a new payment record
func CreatePayment(req PaymentRequest, transactionID string) (*Payment, error) {
	paymentID := generateID()
	cardLastFour := req.CardNumber[len(req.CardNumber)-4:]

	query := `
		INSERT INTO payments (id, order_id, amount, currency, status, card_last_four, transaction_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, order_id, amount, currency, status, card_last_four, transaction_id, created_at, updated_at`

	var payment Payment
	err := DB.QueryRow(query, paymentID, req.OrderID, req.Amount, req.Currency, "completed", cardLastFour, transactionID).Scan(
		&payment.ID, &payment.OrderID, &payment.Amount, &payment.Currency,
		&payment.Status, &payment.CardLastFour, &payment.TransactionID,
		&payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// GetPayment retrieves a payment by its ID
func GetPayment(paymentID string) (*Payment, error) {
	query := `SELECT id, order_id, amount, currency, status, card_last_four, transaction_id, created_at, updated_at FROM payments WHERE id = $1`

	var payment Payment
	err := DB.QueryRow(query, paymentID).Scan(
		&payment.ID, &payment.OrderID, &payment.Amount, &payment.Currency,
		&payment.Status, &payment.CardLastFour, &payment.TransactionID,
		&payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, err
	}

	return &payment, nil
}

// GetPaymentByOrderID retrieves a payment by its order ID
func GetPaymentByOrderID(orderID string) (*Payment, error) {
	query := `SELECT id, order_id, amount, currency, status, card_last_four, transaction_id, created_at, updated_at FROM payments WHERE order_id = $1`

	var payment Payment
	err := DB.QueryRow(query, orderID).Scan(
		&payment.ID, &payment.OrderID, &payment.Amount, &payment.Currency,
		&payment.Status, &payment.CardLastFour, &payment.TransactionID,
		&payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found for this order")
		}
		return nil, err
	}

	return &payment, nil
}

// UpdatePaymentStatus updates the status of a payment
func UpdatePaymentStatus(paymentID, status string) error {
	query := `UPDATE payments SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := DB.Exec(query, status, paymentID)
	return err
}

// ValidateCardNumber validates the card number format
func ValidateCardNumber(cardNumber string) bool {
	// Remove spaces and dashes
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")

	// Check if it's all digits and has valid length (13-19 digits)
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}

	for _, r := range cardNumber {
		if r < '0' || r > '9' {
			return false
		}
	}

	// Basic Luhn algorithm check
	return luhnCheck(cardNumber)
}

// ValidateExpiryDate validates the expiry date format (MM/YY or MM/YYYY)
func ValidateExpiryDate(expiryDate string) bool {
	// Expected format: MM/YY or MM/YYYY
	parts := strings.Split(expiryDate, "/")
	if len(parts) != 2 {
		return false
	}

	month := strings.TrimSpace(parts[0])
	year := strings.TrimSpace(parts[1])

	// Validate month (01-12)
	if len(month) != 2 {
		return false
	}

	// Validate year (2 or 4 digits)
	if len(year) != 2 && len(year) != 4 {
		return false
	}

	// Check if all characters are digits
	for _, r := range month {
		if r < '0' || r > '9' {
			return false
		}
	}

	for _, r := range year {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

// ValidateCVV validates the CVV format (3-4 digits)
func ValidateCVV(cvv string) bool {
	// CVV should be 3 or 4 digits
	if len(cvv) < 3 || len(cvv) > 4 {
		return false
	}

	for _, r := range cvv {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

// luhnCheck implements the Luhn algorithm to validate card numbers
func luhnCheck(cardNumber string) bool {
	nDigits := len(cardNumber)
	parity := nDigits % 2
	total := 0

	for i, r := range cardNumber {
		digit := int(r - '0')
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
		 }
		}
		total += digit
	}

	return total%10 == 0
}

// generateID generates a unique ID
func generateID() string {
	// In a real application, you'd want to use a proper UUID generator
	// For now, we'll use a simple timestamp-based ID
	return fmt.Sprintf("pay_%d", time.Now().UnixNano())
}