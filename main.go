package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func LogInfo(msg string) {
	log.Println("INFO - ", msg)
}

func LogWarning(msg string) {
	log.Println("WARNING - ", msg)
}

func LogError(msg string) {
	log.Println("ERROR - ", msg)
}

type Product struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

var products []Product

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(products) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No products found"})
		return
	}

	json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range products {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		LogError("Failed to decode JSON: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	product.ID = strconv.Itoa(rand.Intn(10000000))
	products = append(products, product)

	response := map[string]interface{}{
		"message": "Product added successfully",
		"id":      product.ID,
	}
	json.NewEncoder(w).Encode(response)

}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var updatedProduct Product

	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		LogError("Failed to decode JSON: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	for index, product := range products {
		if product.ID == params["id"] {

			products[index].Name = updatedProduct.Name
			products[index].Price = updatedProduct.Price
			products[index].Quantity = updatedProduct.Quantity

			response := map[string]interface{}{
				"id":      params["id"],
				"message": "Product updated successfully",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range products {
		if item.ID == params["id"] {
			products = append(products[:index], products[index+1:]...)

			response := map[string]interface{}{
				"id":      params["id"],
				"message": "Product deleted successfully",
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/products", addProduct).Methods("POST")
	r.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")
	log.Println("Server started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
