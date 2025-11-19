package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CartItem struct {
	ProductID     string  `json:"productId"`
	ProductName   string  `json:"productName"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	SelectedSize  string  `json:"selectedSize"`
	SelectedColor string  `json:"selectedColor"`
}

type Cart struct {
	ID        string      `json:"id"`
	UserID    string      `json:"userId"`
	Items     []CartItem  `json:"items"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"userId"`
	Items     []CartItem  `json:"items"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

var carts = make(map[string]Cart)
var orders = make(map[string]Order)
var orderCounter = 0

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
	cartID := fmt.Sprintf("cart_%d", time.Now().UnixNano())

	cart := Cart{
		ID:        cartID,
		UserID:    userID,
		Items:     []CartItem{},
		Total:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	carts[cartID] = cart
	json.NewEncoder(w).Encode(cart)
}

func addItemToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	cartID := vars["cartId"]

	var item CartItem
	json.NewDecoder(r.Body).Decode(&item)

	cart, exists := carts[cartID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cart not found"})
		return
	}

	cart.Items = append(cart.Items, item)
	cart.Total += item.Price * float64(item.Quantity)
	cart.UpdatedAt = time.Now()

	carts[cartID] = cart
	json.NewEncoder(w).Encode(cart)
}

func removeItemFromCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	cartID := vars["cartId"]
	productID := vars["productId"]

	cart, exists := carts[cartID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cart not found"})
		return
	}

	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Total -= item.Price * float64(item.Quantity)
			cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			break
		}
	}

	cart.UpdatedAt = time.Now()
	carts[cartID] = cart
	json.NewEncoder(w).Encode(cart)
}

func getCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	cartID := vars["cartId"]

	cart, exists := carts[cartID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cart not found"})
		return
	}

	json.NewEncoder(w).Encode(cart)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	cartID := vars["cartId"]

	cart, exists := carts[cartID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cart not found"})
		return
	}

	orderCounter++
	orderID := fmt.Sprintf("ORD_%d", orderCounter)

	order := Order{
		ID:        orderID,
		UserID:    cart.UserID,
		Items:     cart.Items,
		Total:     cart.Total,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	orders[orderID] = order

	delete(carts, cartID)

	json.NewEncoder(w).Encode(order)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderID := vars["orderId"]

	order, exists := orders[orderID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
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

	order, exists := orders[orderID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
		return
	}

	order.Status = req["status"]
	order.UpdatedAt = time.Now()
	orders[orderID] = order

	json.NewEncoder(w).Encode(order)
}

func getUserOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userID := vars["userId"]

	var userOrders []Order
	for _, order := range orders {
		if order.UserID == userID {
			userOrders = append(userOrders, order)
		}
	}

	json.NewEncoder(w).Encode(userOrders)
}

func main() {
	r := mux.NewRouter()

	r.Use(enableCORS)

	r.HandleFunc("/api/carts", createCart).Methods("POST")
	r.HandleFunc("/api/carts/{cartId}", getCart).Methods("GET")
	r.HandleFunc("/api/carts/{cartId}/items", addItemToCart).Methods("POST")
	r.HandleFunc("/api/carts/{cartId}/items/{productId}", removeItemFromCart).Methods("DELETE")

	r.HandleFunc("/api/orders", createOrder).Methods("POST")
	r.HandleFunc("/api/orders/{orderId}", getOrder).Methods("GET")
	r.HandleFunc("/api/orders/{orderId}/status", updateOrderStatus).Methods("PUT")
	r.HandleFunc("/api/users/{userId}/orders", getUserOrders).Methods("GET")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	port := ":8002"
	fmt.Printf("Cart & Order Service running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
