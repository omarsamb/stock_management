package main

import (
	"log"
	"stock_management/config"
	"stock_management/db"
	"stock_management/routes"
	"stock_management/services"
)

func main() {
	// Load Configuration
	appConfig := config.LoadConfig() // Call LoadConfig to get configuration

	// Initialize Database
	gormDB, err := db.InitDatabase(appConfig.DatabasePath) // Use config value
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialiser les services
	servicesManager := services.InitServices(gormDB, appConfig.JWTSecret)

	// Configurer les routes
	router := routes.SetupRoutes(servicesManager)

	// Démarrer le serveur
	port := appConfig.ServerPort
	if port == "" {
		port = "8080"
	}

	log.Printf("Serveur démarré sur le port %s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Échec du démarrage du serveur: %v", err)
	}

}
