package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var DB *sql.DB

func getDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   getEnvOrDefault("DB_NAME", "tournyaka"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ConnectDatabase() {
	config := getDBConfig()

	// Use proper string escaping for connection parameters
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password='%s' dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Test connection immediately
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Configure connection pool
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to database")
}

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}