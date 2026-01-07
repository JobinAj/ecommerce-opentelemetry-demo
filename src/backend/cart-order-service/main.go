package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"cart-order-service/db"
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

	var req map[string]string
	json.NewDecoder(r.Body).Decode(&req)

	userID := req["userId"]

	cart, err := db.CreateCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cart)
}

func addItemToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	cartID := vars["cartId"]

	var item db.CartItem
	json.NewDecoder(r.Body).Decode(&item)

	err := db.AddItemToCart(cartID, item)
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

func main() {
	db.InitDB()
	defer db.CloseDB()

	r := mux.NewRouter()

	r.Use(enableCORS)

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

	port := ":8002"
	fmt.Printf("Cart & Order Service running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
