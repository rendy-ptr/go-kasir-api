package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{
		ID:          1,
		Name:        "Snacks",
		Description: "Chips, biscuits, crackers, and other light snacks",
	},
	{
		ID:          2,
		Name:        "Instant Food",
		Description: "Instant noodles, canned food, and ready-to-eat meals",
	},
	{
		ID:          3,
		Name:        "Frozen Food",
		Description: "Frozen meat, nuggets, sausages, and frozen vegetables",
	},
	{
		ID:          4,
		Name:        "Beverages",
		Description: "Soft drinks, mineral water, juice, and tea",
	},
	{
		ID:          5,
		Name:        "Dairy Products",
		Description: "Milk, cheese, yogurt, and other dairy-based products",
	},
	{
		ID:          6,
		Name:        "Bakery",
		Description: "Bread, cakes, pastries, and baked goods",
	},
	{
		ID:          7,
		Name:        "Condiments",
		Description: "Sauces, ketchup, chili sauce, soy sauce, and seasonings",
	},
	{
		ID:          8,
		Name:        "Traditional Food",
		Description: "Local and traditional food products",
	},
}

func getCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for _, category := range categories {
		if category.ID == id {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(category); err != nil {
				log.Println("Error encoding JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "Category not existing", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for index, category := range categories {
		if category.ID == id {
			categories = append(categories[:index], categories[index+1:]...)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]string{
				"message": "Successfully deleted category",
			}); err != nil {
				log.Println("Error encoding JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Category not existing", http.StatusNotFound)
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		if err := json.NewEncoder(w).Encode(categories); err != nil {
			log.Println("Error encoding JSON:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		var newCategory Category
		if err := json.NewDecoder(r.Body).Decode(&newCategory); err != nil {
			log.Println("Error decoding JSON:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newCategory.ID = len(categories) + 1
		categories = append(categories, newCategory)
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newCategory); err != nil {
			log.Println("Error encoding JSON:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}



func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}
	var updateCategory Category
	if err := json.NewDecoder(r.Body).Decode(&updateCategory); err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for index := range categories {
		if categories[index].ID == id {
			updateCategory.ID = id
			categories[index] = updateCategory
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(updateCategory); err != nil {
				log.Println("Error encoding JSON:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Category not existing", http.StatusNotFound)
}

func categoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		getCategoryById(w, r)

	case http.MethodPut:
		updateCategory(w, r)

	case http.MethodDelete:
		deleteCategory(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/api/categories", categoryHandler)
	http.HandleFunc("/api/categories/", categoryByIDHandler)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
