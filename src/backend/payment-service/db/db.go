package db

import (
	"crypto/md5"
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

// Payment represents a payment in the system
type Payment struct {
	ID            string    `json:"id"`
	OrderID       string    `json:"orderId"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	CardLastFour  string    `json:"cardLastFour"`
	TransactionID string    `json:"transactionId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID    string  `json:"orderId"`
	Amount     float64 `json:"amount"`
	CardNumber string  `json:"cardNumber"`
	CardHolder string  `json:"cardHolder"`
	ExpiryDate string  `json:"expiryDate"`
	CVV        string  `json:"cvv"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Payment Payment `json:"payment,omitempty"`
}

// CreatePayment creates a new payment in the database
func CreatePayment(req PaymentRequest) (*Payment, error) {
	paymentID := fmt.Sprintf("PAY_%d", time.Now().UnixNano())
	transactionID := generateTransactionID()
	lastFour := req.CardNumber[len(req.CardNumber)-4:]
	now := time.Now()

	// Insert the payment
	_, err := DB.Exec(`
		INSERT INTO payments (id, order_id, amount, currency, status, card_last_four, transaction_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		paymentID, req.OrderID, req.Amount, "USD", "completed", lastFour, transactionID, now, now)
	if err != nil {
		return nil, err
	}

	payment := &Payment{
		ID:            paymentID,
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		Currency:      "USD",
		Status:        "completed",
		CardLastFour:  lastFour,
		TransactionID: transactionID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	return payment, nil
}

// GetPayment retrieves a payment by its ID
func GetPayment(paymentID string) (*Payment, error) {
	var payment Payment
	var amount sql.NullFloat64
	var currency, status, cardLastFour, transactionID sql.NullString
	var createdAt, updatedAt time.Time

	err := DB.QueryRow(`
		SELECT id, order_id, amount, currency, status, card_last_four, transaction_id, created_at, updated_at
		FROM payments 
		WHERE id = $1`, paymentID).Scan(
		&payment.ID,
		&payment.OrderID,
		&amount,
		&currency,
		&status,
		&cardLastFour,
		&transactionID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, err
	}

	payment.Amount = amount.Float64
	if currency.Valid {
		payment.Currency = currency.String
	} else {
		payment.Currency = "USD"
	}
	payment.Status = status.String
	payment.CardLastFour = cardLastFour.String
	payment.TransactionID = transactionID.String
	payment.CreatedAt = createdAt
	payment.UpdatedAt = updatedAt

	return &payment, nil
}

// GetPaymentByOrderID retrieves a payment by its order ID
func GetPaymentByOrderID(orderID string) (*Payment, error) {
	var payment Payment
	var amount sql.NullFloat64
	var currency, status, cardLastFour, transactionID sql.NullString
	var createdAt, updatedAt time.Time

	err := DB.QueryRow(`
		SELECT id, order_id, amount, currency, status, card_last_four, transaction_id, created_at, updated_at
		FROM payments 
		WHERE order_id = $1`, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&amount,
		&currency,
		&status,
		&cardLastFour,
		&transactionID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found for this order")
		}
		return nil, err
	}

	payment.Amount = amount.Float64
	if currency.Valid {
		payment.Currency = currency.String
	} else {
		payment.Currency = "USD"
	}
	payment.Status = status.String
	payment.CardLastFour = cardLastFour.String
	payment.TransactionID = transactionID.String
	payment.CreatedAt = createdAt
	payment.UpdatedAt = updatedAt

	return &payment, nil
}

// UpdatePaymentStatus updates the status of a payment
func UpdatePaymentStatus(paymentID, status string) error {
	_, err := DB.Exec(`
		UPDATE payments 
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`, status, paymentID)
	return err
}

// ValidateCardNumber validates a card number using Luhn algorithm
func ValidateCardNumber(cardNumber string) bool {
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}

	sum := 0
	isEven := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		if isEven {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isEven = !isEven
	}

	return sum%10 == 0
}

// ValidateExpiryDate validates the expiry date
func ValidateExpiryDate(expiryDate string) bool {
	if len(expiryDate) != 5 || expiryDate[2] != '/' {
		return false
	}

	month := expiryDate[:2]
	year := expiryDate[3:]

	var m, y int
	fmt.Sscanf(month, "%d", &m)
	fmt.Sscanf(year, "%d", &y)

	if m < 1 || m > 12 {
		return false
	}

	currentYear := time.Now().Year() % 100
	if y < currentYear {
		return false
	}

	return true
}

// ValidateCVV validates the CVV
func ValidateCVV(cvv string) bool {
	return len(cvv) >= 3 && len(cvv) <= 4
}

// generateTransactionID generates a unique transaction ID
func generateTransactionID() string {
	data := []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)[:16]
}