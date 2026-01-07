package routes

import (
	"stock_management/handlers"
	"stock_management/middleware"
	"stock_management/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(sm *services.ServicesManager) *gin.Engine {
	router := gin.Default()

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(sm.AccountService, sm.WhatsAppService, sm.JWTSecret)
	articleHandler := handlers.NewArticleHandler(sm.ArticleService)
	stockHandler := handlers.NewStockHandler(sm.StockService)
	subscriptionHandler := handlers.NewSubscriptionHandler(sm.SubscriptionService, sm.AccountService)
	shopHandler := handlers.NewShopHandler(sm.ShopService)
	transferHandler := handlers.NewTransferHandler(sm.StockService)
	dashboardHandler := handlers.NewDashboardHandler(sm.StockService)

	// API Routes
	api := router.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Auth Public Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/verify", authHandler.VerifyPhone)
		}

		// Webhooks (Public)
		api.POST("/subscription/webhook/paydunya", subscriptionHandler.HandlePayDunyaWebhook)

		// Protected Routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(sm.JWTSecret))
		protected.Use(middleware.EnforceSubscriptionMiddleware(sm.DB))
		{
			// Shops
			protected.POST("/shops", shopHandler.CreateShop)
			protected.GET("/shops", shopHandler.ListShops)

			// Users
			protected.POST("/users/invite", authHandler.InviteUser)
			protected.GET("/users", authHandler.ListUsers)
			protected.PUT("/users/:id", authHandler.UpdateUser)
			protected.DELETE("/users/:id", authHandler.DeleteUser)

			// Subscription
			protected.POST("/subscription/select", subscriptionHandler.SelectPlan)

			// Profile
			protected.PUT("/auth/profile", authHandler.UpdateProfile)
			protected.POST("/auth/change-password", authHandler.ChangePassword)
			protected.PUT("/auth/theme", authHandler.UpdateTheme)

			// Articles
			protected.POST("/articles", articleHandler.CreateArticle)
			protected.GET("/articles", articleHandler.ListArticles)
			protected.PUT("/articles/:id", articleHandler.UpdateArticle)
			protected.POST("/articles/import", articleHandler.ImportArticles)

			// Dashboard
			protected.GET("/dashboard/stats", dashboardHandler.GetStats)

			// Stocks
			protected.POST("/stocks/movement", stockHandler.RecordMovement)
			protected.GET("/stocks/levels", stockHandler.ListStockLevels)
			protected.GET("/stocks/movements", stockHandler.ListMovements)

			// Transfers
			protected.POST("/transfers", transferHandler.InitiateTransfer)
			protected.POST("/transfers/:id/receive", transferHandler.ReceiveTransfer)
			protected.GET("/transfers", transferHandler.ListTransfers)
		}
	}

	return router
}
