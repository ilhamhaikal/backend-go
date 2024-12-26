package routes

import (
	"encoding/json"
	"net/http"
	"tournyaka-backend/controllers"
	"tournyaka-backend/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"version": "1.0.0",
		})
	}).Methods("GET")

	// Root handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":       "Welcome to Tournyaka API",
			"documentation": "/api/docs",
			"version":       "1.0.0",
		})
	}).Methods("GET")

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes
	api.HandleFunc("/login", controllers.Login).Methods("POST")
	api.HandleFunc("/register", controllers.Register).Methods("POST")
	api.HandleFunc("/refresh-token", controllers.RefreshToken).Methods("POST")

	// Protected routes
	protected := api.PathPrefix("/user").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/profile", controllers.GetUserProfile).Methods("GET")
	protected.HandleFunc("/update", controllers.UpdateUser).Methods("PUT")

	return router
}
