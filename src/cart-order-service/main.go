package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cart-order-service/db"
	"cart-order-service/telemetry"

	"github.com/gorilla/mux"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

func createCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Feature Flag Check
	client := openfeature.NewClient("cart-order-service")
	failureEnabled, err := client.BooleanValue(context.Background(), "cartServiceFailure", false, openfeature.EvaluationContext{})
	if err == nil && failureEnabled {
		log.Println("Simulated Cart Service Failure triggered")
		http.Error(w, "Simulated Cart Service Failure", http.StatusInternalServerError)
		return
	}

	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateCart: Invalid request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID := req["userId"]
	log.Printf("CreateCart request for userID: %s", userID)

	cart, err := db.CreateCart(userID)
	if err != nil {
		log.Printf("CreateCart: DB error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("CreateCart successful. CartID: %s", cart.ID)
	json.NewEncoder(w).Encode(cart)
}

func addItemToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	cartID := vars["cartId"]

	var item db.CartItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		log.Printf("AddItemToCart: Invalid request body for cart %s: %v", cartID, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("AddItemToCart request for cart %s. ProductID: %s, Quantity: %d", cartID, item.ProductID, item.Quantity)

	err := db.AddItemToCart(cartID, item)
	if err != nil {
		log.Printf("AddItemToCart: DB error for cart %s: %v", cartID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cart, err := db.GetCart(cartID)
	if err != nil {
		log.Printf("AddItemToCart: Failed to retrieve cart %s after adding item: %v", cartID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("AddItemToCart successful for cart %s", cartID)
	json.NewEncoder(w).Encode(cart)
}

func removeItemFromCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	cartID := vars["cartId"]
	productID := vars["productId"]

	err := db.RemoveItemFromCart(cartID, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cart, err := db.GetCart(cartID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cart)
}

func getCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	cartID := vars["cartId"]

	cart, err := db.GetCart(cartID)
	if err != nil {
		if err.Error() == "cart not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Cart not found"})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cart)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	cartID := vars["cartId"]

	order, err := db.CreateOrder(cartID)
	if err != nil {
		if err.Error() == "cart not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Cart not found"})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderID := vars["orderId"]

	order, err := db.GetOrder(orderID)
	if err != nil {
		if err.Error() == "order not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func updateOrderStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderID := vars["orderId"]

	var req map[string]string
	json.NewDecoder(r.Body).Decode(&req)

	err := db.UpdateOrderStatus(orderID, req["status"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order, err := db.GetOrder(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func getUserOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID := vars["userId"]

	orders, err := db.GetUserOrders(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Signup: Invalid request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Signup payload received - Email: '%s', Name: '%s', Password Length: %d", req.Email, req.Name, len(req.Password))
	if len(req.Password) > 0 {
		log.Printf("Signup password first char: '%c'", req.Password[0])
	} else {
		log.Println("Signup password is empty")
	}

	log.Printf("Signup attempt for email: %s", req.Email)

	// Check if user exists
	existingUser, _ := db.GetUserByEmail(req.Email)
	if existingUser != nil {
		log.Printf("Signup: User already exists: %s", req.Email)
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Create new user
	// Generate ID
	userID := fmt.Sprintf("user_%d", time.Now().UnixNano())

	user := db.User{
		ID:           userID,
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: req.Password, // Insecure demo only
	}

	if err := db.CreateUser(user); err != nil {
		log.Printf("Signup: Failed to create user %s: %v", req.Email, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Signup successful for user: %s", req.Email)
	json.NewEncoder(w).Encode(user)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Login: Invalid request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Login attempt for email: %s", req.Email)

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("Login: User not found or DB error for email %s: %v", req.Email, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		return
	}

	// Simple password check (Insecure demo only)
	if user.PasswordHash != req.Password {
		log.Printf("Login: Password mismatch for email %s. Stored: '%s', Provided: '%s'", req.Email, user.PasswordHash, req.Password)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		return
	}

	log.Printf("Login successful for email: %s", req.Email)
	json.NewEncoder(w).Encode(user)
}

func main() {
	// Initialize OpenTelemetry
	ctx := context.Background()
	shutdownTracer := telemetry.InitTracer(ctx)
	defer shutdownTracer(ctx)

	// Initialize OpenFeature provider
	provider, err := flagd.NewProvider(
		flagd.WithHost("otel-flagd.apps.svc.cluster.local"),
		flagd.WithPort(8013),
	)
	if err != nil {
		log.Printf("Warning: Failed to create flagd provider: %v", err)
	} else {
		openfeature.SetProvider(provider)
	}

	db.InitDB()
	defer db.CloseDB()

	r := mux.NewRouter()

	r.Use(enableCORS)

	// Auth Routes
	r.HandleFunc("/api/signup", handleSignup).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/login", handleLogin).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/carts", createCart).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/carts/{cartId}", getCart).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/carts/{cartId}/items", addItemToCart).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/carts/{cartId}/items/{productId}", removeItemFromCart).Methods("DELETE", "OPTIONS")

	r.HandleFunc("/api/carts/{cartId}/orders", createOrder).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/orders/{orderId}", getOrder).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/orders/{orderId}/status", updateOrderStatus).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/users/{userId}/orders", getUserOrders).Methods("GET", "OPTIONS")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET", "OPTIONS")

	port := "8002"
	fmt.Printf("Cart & Order Service running on http://localhost:%s\n", port)
	handler := otelhttp.NewHandler(r, "cart-order-service")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}
