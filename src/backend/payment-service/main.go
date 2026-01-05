package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"versace-payment-service/db"
)

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

	payment, err := db.CreatePayment(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(db.PaymentResponse{
		Success: true,
		Message: "Payment processed successfully",
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payment, err := db.GetPayment(paymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(payment)
}

func main() {
	db.InitDB()
	defer db.CloseDB()

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
