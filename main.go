package main

import (
	"encoding/json"
	"fmt"
	"kasirapi/database"
	"kasirapi/handlers"
	"kasirapi/repositories"
	"kasirapi/services"
	"log"
	"net/http"

	"github.com/lynxsecurity/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	log.Println("PORT:", config.Port)

	log.Println("DB_CONN:", config.DBConn)

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	//GET localhost:8088/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup routes
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	//Endpoint localhost:8080/api/categories
	//http.HandleFunc("/api/categories", getpostCategory)

	//localhost:8088//health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})
	fmt.Println("server running di localhost:" + config.Port)

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("gagal running")
	}
}
