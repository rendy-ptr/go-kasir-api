package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int     `json:"id"`
	Nama  string  `json:"nama"`
	Harga float64 `json:"harga"`
	Stok  int     `json:"stok"`
}

var produk = []Produk{
	{
		ID:    1,
		Nama:  "Indomie Godog",
		Harga: 3500,
		Stok:  10,
	},
	{
		ID:    2,
		Nama:  "Vit 1000 ml",
		Harga: 3000,
		Stok:  40,
	},
	{
		ID:    3,
		Nama:  "Kecap",
		Harga: 12000,
		Stok:  20,
	},
}

func getProdukById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(p); err != nil {
				log.Println("Error encoding JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{
				"message": "Sukses Delete Produk",
			}); err != nil {
				log.Println("Error encoding JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":  "OK",
		"message": "Api is running",
	}); err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func produkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(produk); err != nil {
			log.Println("Error encoding JSON:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		var produkBaru Produk
		if err := json.NewDecoder(r.Body).Decode(&produkBaru); err != nil {
			log.Println("Error decoding JSON:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		produkBaru.ID = len(produk) + 1
		produk = append(produk, produkBaru)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(produkBaru); err != nil {
			log.Println("Error encoding JSON:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func produkByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		getProdukById(w, r)

	case http.MethodPut:
		updateProduk(w, r)

	case http.MethodDelete:
		deleteProduk(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	var updateProduk Produk
	if err := json.NewDecoder(r.Body).Decode(&updateProduk); err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(updateProduk); err != nil {
				log.Println("Error encoding JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/produk", produkHandler)
	http.HandleFunc("/api/produk/", produkByIDHandler)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
