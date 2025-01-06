package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"shelter-it-be/internal/config"
	"shelter-it-be/internal/database"
	"shelter-it-be/internal/handler"
	"shelter-it-be/internal/middleware"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Init DB
	db, err := database.NewDB(cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect DB ", err)
	}
	defer db.DB.Close()

	// Set Gin mode
	if cfg.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Init router with middleware
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Init handlers with JWT config
	authHandler := handler.NewAuthHandler(db, []byte(cfg.JWT.Secret))

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware([]byte(cfg.JWT.Secret)))
	{
		protected.POST("/refresh-token", authHandler.RefreshToken)
		protected.POST("/logout", authHandler.Logout)
		protected.GET("/profile", getUserProfile)
	}

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server is running on %s", serverAddr)

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Server failed to start: ", err)
	}
}

func getUserProfile(c *gin.Context) {
	uid, _ := c.Get("uid")
	email, _ := c.Get("email")
	c.JSON(http.StatusOK, gin.H{
		"uid":   uid,
		"email": email,
	})
}
