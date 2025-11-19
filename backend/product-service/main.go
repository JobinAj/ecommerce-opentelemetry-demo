package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Product struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Price    float64  `json:"price"`
	Image    string   `json:"image"`
	Desc     string   `json:"description"`
	Rating   float64  `json:"rating"`
	Reviews  int      `json:"reviews"`
	Sizes    []string `json:"sizes"`
	Colors   []string `json:"colors"`
	InStock  bool     `json:"inStock"`
}

var products = []Product{
	{
		ID:       "1",
		Name:     "Silk Baroque Shirt",
		Category: "Tops",
		Price:    1200,
		Image:    "https://images.pexels.com/photos/3622613/pexels-photo-3622613.jpeg",
		Desc:     "Luxurious silk shirt with signature Versace Baroque pattern",
		Rating:   4.8,
		Reviews:  124,
		Sizes:    []string{"XS", "S", "M", "L", "XL"},
		Colors:   []string{"Black", "Gold", "White"},
		InStock:  true,
	},
	{
		ID:       "2",
		Name:     "Gold Medusa Blazer",
		Category: "Outerwear",
		Price:    2500,
		Image:    "https://images.pexels.com/photos/3622614/pexels-photo-3622614.jpeg",
		Desc:     "Statement blazer featuring iconic Medusa head emblem",
		Rating:   4.9,
		Reviews:  89,
		Sizes:    []string{"XS", "S", "M", "L", "XL"},
		Colors:   []string{"Black", "Navy", "Charcoal"},
		InStock:  true,
	},
	{
		ID:       "3",
		Name:     "Versace Print T-Shirt",
		Category: "Tops",
		Price:    450,
		Image:    "https://images.pexels.com/photos/3622615/pexels-photo-3622615.jpeg",
		Desc:     "Classic cotton t-shirt with Versace logo print",
		Rating:   4.6,
		Reviews:  256,
		Sizes:    []string{"XS", "S", "M", "L", "XL", "XXL"},
		Colors:   []string{"White", "Black", "Red", "Navy"},
		InStock:  true,
	},
	{
		ID:       "4",
		Name:     "Tailored Silk Trousers",
		Category: "Bottoms",
		Price:    1800,
		Image:    "https://images.pexels.com/photos/3622616/pexels-photo-3622616.jpeg",
		Desc:     "Elegant silk trousers with perfect drape",
		Rating:   4.7,
		Reviews:  142,
		Sizes:    []string{"XS", "S", "M", "L", "XL"},
		Colors:   []string{"Black", "White", "Beige"},
		InStock:  true,
	},
	{
		ID:       "5",
		Name:     "Black Leather Jacket",
		Category: "Outerwear",
		Price:    3200,
		Image:    "https://images.pexels.com/photos/3622617/pexels-photo-3622617.jpeg",
		Desc:     "Premium leather jacket with signature detailing",
		Rating:   4.9,
		Reviews:  198,
		Sizes:    []string{"XS", "S", "M", "L", "XL"},
		Colors:   []string{"Black", "Brown"},
		InStock:  true,
	},
	{
		ID:       "6",
		Name:     "Gold Chain Dress",
		Category: "Dresses",
		Price:    2800,
		Image:    "https://images.pexels.com/photos/3622618/pexels-photo-3622618.jpeg",
		Desc:     "Stunning dress with gold chain embellishments",
		Rating:   4.8,
		Reviews:  167,
		Sizes:    []string{"XS", "S", "M", "L"},
		Colors:   []string{"Black", "Gold", "Silver"},
		InStock:  true,
	},
	{
		ID:       "7",
		Name:     "Premium Denim Jeans",
		Category: "Bottoms",
		Price:    950,
		Image:    "https://images.pexels.com/photos/3622619/pexels-photo-3622619.jpeg",
		Desc:     "High-quality denim with Versace branding",
		Rating:   4.7,
		Reviews:  203,
		Sizes:    []string{"24", "25", "26", "27", "28", "29", "30", "31", "32"},
		Colors:   []string{"Dark Blue", "Light Blue", "Black"},
		InStock:  true,
	},
	{
		ID:       "8",
		Name:     "Silk Evening Gown",
		Category: "Dresses",
		Price:    4500,
		Image:    "https://images.pexels.com/photos/3622620/pexels-photo-3622620.jpeg",
		Desc:     "Breathtaking silk gown for special occasions",
		Rating:   5.0,
		Reviews:  87,
		Sizes:    []string{"XS", "S", "M", "L"},
		Colors:   []string{"Black", "Red", "White"},
		InStock:  true,
	},
}

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
	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")

	var filtered []Product

	for _, p := range products {
		if category != "" && p.Category != category {
			continue
		}
		if search != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(search)) {
			continue
		}
		filtered = append(filtered, p)
	}

	json.NewEncoder(w).Encode(filtered)
}

func getProductByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	for _, p := range products {
		if p.ID == id {
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categoryMap := make(map[string]bool)
	var categories []string

	for _, p := range products {
		if !categoryMap[p.Category] {
			categoryMap[p.Category] = true
			categories = append(categories, p.Category)
		}
	}

	json.NewEncoder(w).Encode(categories)
}

func updateStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req map[string]interface{}
	json.NewDecoder(r.Body).Decode(&req)

	productID := req["productId"].(string)
	quantity := int(req["quantity"].(float64))

	for i, p := range products {
		if p.ID == productID {
			if p.InStock && quantity > 0 {
				products[i].InStock = true
				json.NewEncoder(w).Encode(map[string]bool{"success": true})
				return
			}
		}
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": "Stock update failed"})
}

func searchProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("q")
	minPrice := r.URL.Query().Get("minPrice")
	maxPrice := r.URL.Query().Get("maxPrice")

	min := 0.0
	max := 100000.0

	if minPrice != "" {
		if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
			min = val
		}
	}
	if maxPrice != "" {
		if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			max = val
		}
	}

	var results []Product

	for _, p := range products {
		if (strings.Contains(strings.ToLower(p.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(p.Desc), strings.ToLower(query))) &&
			p.Price >= min && p.Price <= max {
			results = append(results, p)
		}
	}

	json.NewEncoder(w).Encode(results)
}

func main() {
	r := mux.NewRouter()

	r.Use(enableCORS)

	r.HandleFunc("/api/products", getAllProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProductByID).Methods("GET")
	r.HandleFunc("/api/categories", getCategories).Methods("GET")
	r.HandleFunc("/api/search", searchProducts).Methods("GET")
	r.HandleFunc("/api/stock/update", updateStock).Methods("POST")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	port := ":8001"
	fmt.Printf("Product Service running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
