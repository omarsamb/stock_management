package main

import (
	"log"
	"net/http"
	"stock_management/backend/database"
	"stock_management/backend/handlers"
	"stock_management/backend/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)

	// Create router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Authentication routes (no auth required)
	api.HandleFunc("/signup", authHandler.Signup).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// Protected routes (auth required)
	api.HandleFunc("/profile", middleware.JWTMiddleware(authHandler.Profile)).Methods("GET")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("Backend server starting on port 8080")
	log.Println("Available endpoints:")
	log.Println("  POST /api/signup   - Create new tenant and admin user")
	log.Println("  POST /api/login    - Authenticate user")
	log.Println("  POST /api/logout   - Logout user")
	log.Println("  GET  /api/profile  - Get user profile (requires auth)")
	log.Println("  GET  /health       - Health check")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
