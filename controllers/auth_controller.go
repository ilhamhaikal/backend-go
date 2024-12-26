package controllers

import (
	"encoding/json"
	"net/http"
	"tournyaka-backend/config"
	"tournyaka-backend/utils"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	var userID int64
	err = config.DB.QueryRow(`
        INSERT INTO users (username, email, password) 
        VALUES ($1, $2, $3)
        RETURNING id`,
		req.Username, req.Email, string(hashedPassword)).Scan(&userID)

	if err != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user_id": userID,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Query untuk mendapatkan data dari kedua tabel
	query := `
        SELECT u.id, u.password, mu.nama, mr.role_name 
        FROM users u 
        LEFT JOIN m_users mu ON u.id = mu.user_id 
        LEFT JOIN m_role mr ON mu.role_id = mr.id 
        WHERE u.email = $1
    `

	var user struct {
		ID       uint
		Password string
		Nama     string
		Role     string
	}

	err := config.DB.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Password, &user.Nama, &user.Role)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"token":  token,
		"user": map[string]interface{}{
			"id":   user.ID,
			"nama": user.Nama,
			"role": user.Role,
		},
	})
}

// GetUserProfile - handler untuk mendapatkan profil pengguna
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Profile retrieved successfully",
		// TODO: Tambahkan data profil user
	})
}

// UpdateUser - handler untuk memperbarui data pengguna
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "User updated successfully",
		// TODO: Tambahkan logika update user
	})
}

// RefreshToken - handler untuk memperbarui token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Token refreshed successfully",
		// TODO: Tambahkan logika refresh token
	})
}
