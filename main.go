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
)

func main() {
	// init DB
	config.InitDB()

	// gin
	r := gin.Default()

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
		&models.CartItem{},
	)

	// routes
	router.Setup(r)

	r.Run() // :8080 by default
}
