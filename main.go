// @title Bogbon API
// @version 1.0
// @description API for the Bogbon e-commerce site (categories, products, cart, and orders).

// @host localhost:8080
// @BasePath /api

package main

import (
	"log"
	"os"

	"bogbon-api/config"
	"bogbon-api/models"
	"bogbon-api/router"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	_ "bogbon-api/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func main() {
	// init DB
	config.InitDB()

	// gin
	r := gin.Default()

	// Serve uploaded images
	r.Static("/uploads", "./uploads")

	// Limit file size
	r.MaxMultipartMemory = 8 << 20 // 8 MB

	// session middleware (cookie store)
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		log.Fatal("SESSION_SECRET must be set")
	}
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400 * 30, // 30 days
	})
	r.Use(sessions.Sessions("bogbon_session", store))

	// auto‑migrate
	config.DB.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
	)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// routes
	router.Setup(r)

	r.Run() // :8080 by default
}
