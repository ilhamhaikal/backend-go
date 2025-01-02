package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtKey     = []byte(os.Getenv("JWT_SECRET_KEY"))
	usedTokens = sync.Map{} // Track used tokens
)

func GenerateJWT(userID uint) (string, error) {
	// Periksa jika user sudah memiliki token aktif
	if storedToken, exists := usedTokens.Load(userID); exists {
		if IsTokenValid(storedToken.(string)) {
			return "", fmt.Errorf("user already has an active session")
		}
		// Jika token tidak valid, hapus dari usedTokens
		usedTokens.Delete(userID)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(jwtKey)
	if err == nil {
		usedTokens.Store(userID, tokenString)
	}
	return tokenString, err
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("no token provided")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if token is expired
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token expired")
			}
		}
		return token, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func InvalidateToken(tokenString string) error {
	log.Printf("Attempting to invalidate token: %s", tokenString)
	if tokenString == "" {
		return fmt.Errorf("no token provided")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		log.Printf("Invalid token during logout: %v", err)
		return fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(float64); ok {
			log.Printf("Invalidating token for user ID: %v", uint(userID))
			usedTokens.Delete(uint(userID))
			return nil
		}
	}
	return fmt.Errorf("invalid token claims")
}

func IsTokenValid(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID := uint(claims["user_id"].(float64))
		if storedToken, exists := usedTokens.Load(userID); exists {
			return storedToken == tokenString
		}
	}
	return false
}
