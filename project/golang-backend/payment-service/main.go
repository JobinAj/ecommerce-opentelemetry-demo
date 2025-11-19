package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

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

type PaymentRequest struct {
	OrderID    string  `json:"orderId"`
	Amount     float64 `json:"amount"`
	CardNumber string  `json:"cardNumber"`
	CardHolder string  `json:"cardHolder"`
	ExpiryDate string  `json:"expiryDate"`
	CVV        string  `json:"cvv"`
}

type PaymentResponse struct {
	Success bool    `json:"success"`
	Message string  `json:"message"`
	Payment Payment `json:"payment,omitempty"`
}

var payments = make(map[string]Payment)
var paymentCounter = 0

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateCardNumber(cardNumber string) bool {
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

func validateExpiryDate(expiryDate string) bool {
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

func validateCVV(cvv string) bool {
	return len(cvv) >= 3 && len(cvv) <= 4
}

func generateTransactionID() string {
	data := []byte(fmt.Sprintf("%d", time.Now().UnixNano()))
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)[:16]
}

func processPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	if !validateCardNumber(req.CardNumber) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentResponse{
			Success: false,
			Message: "Invalid card number",
		})
		return
	}

	if !validateExpiryDate(req.ExpiryDate) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentResponse{
			Success: false,
			Message: "Invalid or expired card",
		})
		return
	}

	if !validateCVV(req.CVV) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(PaymentResponse{
			Success: false,
			Message: "Invalid CVV",
		})
		return
	}

	paymentCounter++
	paymentID := fmt.Sprintf("PAY_%d", paymentCounter)
	transactionID := generateTransactionID()
	lastFour := req.CardNumber[len(req.CardNumber)-4:]

	payment := Payment{
		ID:           paymentID,
		OrderID:      req.OrderID,
		Amount:       req.Amount,
		Currency:     "USD",
		Status:       "completed",
		CardLastFour: lastFour,
		TransactionID: transactionID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	payments[paymentID] = payment

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PaymentResponse{
		Success: true,
		Message: "Payment processed successfully",
		Payment: payment,
	})
}

func getPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	paymentID := vars["paymentId"]

	payment, exists := payments[paymentID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payment not found"})
		return
	}

	json.NewEncoder(w).Encode(payment)
}

func getPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderID := vars["orderId"]

	for _, payment := range payments {
		if payment.OrderID == orderID {
			json.NewEncoder(w).Encode(payment)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Payment not found for this order"})
}

func refundPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	paymentID := vars["paymentId"]

	payment, exists := payments[paymentID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payment not found"})
		return
	}

	payment.Status = "refunded"
	payment.UpdatedAt = time.Now()
	payments[paymentID] = payment

	json.NewEncoder(w).Encode(payment)
}

func main() {
	r := mux.NewRouter()

	r.Use(enableCORS)

	r.HandleFunc("/api/payments", processPayment).Methods("POST")
	r.HandleFunc("/api/payments/{paymentId}", getPayment).Methods("GET")
	r.HandleFunc("/api/payments/order/{orderId}", getPaymentByOrderID).Methods("GET")
	r.HandleFunc("/api/payments/{paymentId}/refund", refundPayment).Methods("POST")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	port := ":8003"
	fmt.Printf("Payment Service running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
