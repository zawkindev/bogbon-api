// @title           Bogbon API
// @version         1.0
// @description     Swagger documentation for the Bogbon Gin API.
// @host      localhost:8080
// @BasePath  /api

package main

import (
	"log"
	"os"
	"time"

	"bogbon-api/config"
	"bogbon-api/models"
	"bogbon-api/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	_ "bogbon-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// init DB
	config.InitDB()

	// gin
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://gardening-service.uz"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve uploaded images
	r.Static("/uploads", "./uploads")

	// Limit file size
	// r.MaxMultipartMemory = 8 << 20 // 8 MB

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
		&models.CategoryTranslation{},
		&models.Product{},
		&models.ProductTranslation{},
		&models.ProductImage{},
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
