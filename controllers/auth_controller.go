package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"tournyaka-backend/config"
	"tournyaka-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input
	if !isValidEmail(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}
	if len(req.Username) == 0 {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Insert user ke database
	var userID int
	err = config.DB.QueryRow(
		`INSERT INTO users (username, email, password_hash, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) RETURNING user_id`,
		req.Username, req.Email, string(hashedPassword), time.Now(), time.Now(),
	).Scan(&userID)
	if err != nil {
		log.Println("Error creating user:", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Berikan respons sukses
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Registration successful",
		"user_id": userID,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input
	if !isValidEmail(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Get user from database
	var user struct {
		UserID       int
		PasswordHash string
		Username     string
		Role         int
	}
	err := config.DB.QueryRow(
		`SELECT user_id, password_hash, username, COALESCE(role, 3) as role 
         FROM users 
         WHERE email = $1`,
		req.Email,
	).Scan(&user.UserID, &user.PasswordHash, &user.Username, &user.Role)

	if err != nil {
		log.Println("Error finding user:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Validasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := utils.GenerateJWT(uint(user.UserID))
	if err != nil {
		log.Println("Error generating JWT:", err)
		http.Error(w, "User already logged in", http.StatusForbidden)
		return
	}

	// Return response with role
	redirectPath := "/dashboard" // Default all roles to dashboard
	if user.Role == 3 {          // Only regular users (role 3) go to homepage
		redirectPath = "/"
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"token":  token,
		"user": map[string]interface{}{
			"id":              user.UserID,
			"username":        user.Username,
			"role":            user.Role,
			"redirect":        redirectPath,
			"isAuthenticated": true,
		},
	})
}

// Logout handler
func Logout(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		log.Println("No Authorization header provided")
		http.Error(w, "No token provided", http.StatusBadRequest)
		return
	}

	log.Printf("Authorization header: %s", tokenString)
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	err := utils.InvalidateToken(tokenString)
	if err != nil {
		log.Printf("Error invalidating token: %v", err)
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	log.Println("Logout successful")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Successfully logged out",
	})
}
