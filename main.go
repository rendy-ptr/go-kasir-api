package main

import (
	"fmt"
	"go-kasir-api/database"
	"go-kasir-api/handlers"
	"go-kasir-api/repositories"
	"go-kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// http.HandleFunc("/api/categories", categoryHandler)
	// http.HandleFunc("/api/categories/", categoryByIDHandler)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	http.HandleFunc("/api/category", categoryHandler.HandleCategories)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("API is running")); err != nil {
			log.Println("Error writing response:", err)
		}
	})

	// LOCAL
	// fmt.Println("Server running on http://localhost:8080")
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal("ListenAndServe:", err)
	// }

	// DEPLOYMENT
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }

	addr := "localhost:" + config.Port

	fmt.Println("Server running di http://" + addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("gagal running server", err)
	}

}
