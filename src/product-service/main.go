package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"product-service/db"
	"product-service/telemetry"

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

func getAllProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Feature Flag Check
	client := openfeature.NewClient("product-service")
	failureEnabled, err := client.BooleanValue(context.Background(), "productCatalogFailure", false, openfeature.EvaluationContext{})
	if err == nil && failureEnabled {
		http.Error(w, "Simulated Product Catalog Failure", http.StatusInternalServerError)
		return
	}

	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")

	products, err := db.GetProducts(category, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := db.GetProductByID(id)
	if err != nil {
		if err.Error() == "product not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	categories, err := db.GetCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(categories)
}

func updateStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req map[string]interface{}
	json.NewDecoder(r.Body).Decode(&req)

	productID := req["productId"].(string)
	quantity := int(req["quantity"].(float64))

	err := db.UpdateStock(productID, quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func searchProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")

	products, err := db.SearchProducts(query, minPrice, maxPrice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
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

	r.HandleFunc("/api/products", getAllProducts).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/products/{id}", getProductByID).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/categories", getCategories).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/search", searchProducts).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/stock/update", updateStock).Methods("POST", "OPTIONS")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET", "OPTIONS")

	port := "8001"
	fmt.Printf("Product Service running on http://localhost:%s\n", port)
	handler := otelhttp.NewHandler(r, "product-service")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handler))
}
