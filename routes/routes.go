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

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"version": "1.0.0",
		})
	}).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":       "Welcome to Tournyaka API",
			"documentation": "/api/docs",
			"version":       "1.0.0",
		})
	}).Methods("GET")

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/login", controllers.Login).Methods("POST")
	api.HandleFunc("/register", controllers.Register).Methods("POST")

	protected := api.PathPrefix("/user").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	return router
}
