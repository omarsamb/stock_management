package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	// "strconv" // Not strictly needed for this basic version, but good for future if parsing ints, bools
)

// Config holds application configuration values.
type Config struct {
	JWTSecret    string
	DatabasePath string
	ServerPort   string
}

// AppConfig is a global variable holding the application configuration.
// It will be initialized by LoadConfig().
// While global variables are generally discouraged, they can be acceptable for app-wide configs if managed carefully.
// Alternatively, the config can be passed around explicitly. For this task, we'll use the global.
var AppConfig *Config

// LoadConfig initializes the AppConfig from environment variables or defaults.
// It also returns the loaded config.

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Skipping...")
	}
	// Build the connection string
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback: construire à partir des variables DB_*
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Africa/Dakar",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSLMODE"),
		)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	AppConfig = &Config{
		JWTSecret:    getEnv("JWT_SECRET", "flashcard_secret"),
		DatabasePath: dsn,
		ServerPort:   port,
	}
	return AppConfig
}

// getEnv retrieves an environment variable by key or returns a fallback value if not set.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
