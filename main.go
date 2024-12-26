package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"tournyaka-backend/config"
	"tournyaka-backend/routes"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDatabase()
	defer config.CloseDatabase()

	router := routes.RegisterRoutes()

	// Add CORS middleware
	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	// Add logging
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	// Start server with middleware
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080",
		handlers.CORS(headers, methods, origins)(loggedRouter)))
}
