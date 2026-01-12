package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"payment-service/db"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"

	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

var sc *client.API

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

func processPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Feature Flag Check
	client := openfeature.NewClient("payment-service")
	failureEnabled, err := client.BooleanValue(context.Background(), "paymentServiceFailure", false, openfeature.EvaluationContext{})
	if err == nil && failureEnabled {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(db.PaymentResponse{
			Success: false,
			Message: "Simulated Payment Service Failure",
		})
		return
	}

	var req db.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(db.PaymentResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	if !db.ValidateCardNumber(req.CardNumber) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(db.PaymentResponse{
			Success: false,
			Message: "Invalid card number",
		})
		return
	}

	if !db.ValidateExpiryDate(req.ExpiryDate) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(db.PaymentResponse{
			Success: false,
			Message: "Invalid or expired card",
		})
		return
	}

	if !db.ValidateCVV(req.CVV) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(db.PaymentResponse{
			Success: false,
			Message: "Invalid CVV",
		})
		return
	}

	// Parse expiry date
	var expMonth, expYear string
	if len(req.ExpiryDate) == 5 { // MM/YY
		expMonth = req.ExpiryDate[:2]
		expYear = "20" + req.ExpiryDate[3:]
	} else if len(req.ExpiryDate) == 7 { // MM/YYYY
		expMonth = req.ExpiryDate[:2]
		expYear = req.ExpiryDate[3:]
	} else {
		// Should be caught by validation, but just in case
		expMonth = "12"
		expYear = "2025"
	}

	// 1. Create a Token
	tokenParams := &stripe.TokenParams{
		Card: &stripe.CardParams{
			Number:   stripe.String(req.CardNumber),
			ExpMonth: stripe.String(expMonth),
			ExpYear:  stripe.String(expYear),
			CVC:      stripe.String(req.CVV),
		},
	}

	token, err := sc.Tokens.New(tokenParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(db.PaymentResponse{
			Success: false,
			Message: "Stripe Token Error: " + err.Error(),
		})
		return
	}

	// 2. Create a Charge
	chargeParams := &stripe.ChargeParams{
		Amount:   stripe.Int64(int64(req.Amount * 100)), // Amount in cents
		Currency: stripe.String(req.Currency),
		Source:   &stripe.SourceParams{Token: stripe.String(token.ID)},
	}

	charge, err := sc.Charges.New(chargeParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(db.PaymentResponse{
			Success: false,
			Message: "Stripe Charge Error: " + err.Error(),
		})
		return
	}

	// 3. Save to DB using Stripe Charge ID
	payment, err := db.CreatePayment(req, charge.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(db.PaymentResponse{
			Success: false,
			Message: "Database Error: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(db.PaymentResponse{
		Success: true,
		Message: "Payment processed successfully via Stripe",
		Payment: *payment,
	})
}

func getPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	paymentID := vars["paymentId"]

	payment, err := db.GetPayment(paymentID)
	if err != nil {
		if err.Error() == "payment not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Payment not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(payment)
}

func getPaymentByOrderID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderID := vars["orderId"]

	payment, err := db.GetPaymentByOrderID(orderID)
	if err != nil {
		if err.Error() == "payment not found for this order" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Payment not found for this order"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(payment)
}

func refundPayment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	paymentID := vars["paymentId"]

	err := db.UpdatePaymentStatus(paymentID, "refunded")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	payment, err := db.GetPayment(paymentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(payment)
}

func main() {
	// Initialize OpenFeature provider
	provider, err := flagd.NewProvider()
	if err != nil {
		log.Fatalf("Failed to create flagd provider: %v", err)
	}
	openfeature.SetProvider(provider)

	db.InitDB()
	defer db.CloseDB()

	// Initialize Stripe Client pointing to stripe-mock
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		stripeKey = "sk_test_default_mock_key"
	}
	sc = &client.API{}
	sc.Init(stripeKey, &stripe.Backends{
		API: stripe.GetBackendWithConfig(stripe.APIBackend, &stripe.BackendConfig{
			URL: stripe.String("http://localhost:12111"), // Address of stripe-mock
		}),
	})

	r := mux.NewRouter()

	r.Use(enableCORS)

	r.HandleFunc("/api/payments", processPayment).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/payments/{paymentId}", getPayment).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/payments/order/{orderId}", getPaymentByOrderID).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/payments/{paymentId}/refund", refundPayment).Methods("POST", "OPTIONS")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET", "OPTIONS")

	port := ":8003"
	fmt.Printf("Payment Service running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
