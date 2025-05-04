package main

import (
	"bookstore-api/database"
	"bookstore-api/handlers"
	"bookstore-api/middlewares"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	database.Connect()

	// routes
	r.GET("/books", handlers.GetBooks)
	r.GET("/books/:id", handlers.GetBookByID)
	r.GET("/categories", handlers.GetCategories)
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	auth := r.Group("/")
	auth.Use(middlewares.JWTAuthMiddleware())
	auth.POST("/books", handlers.CreateBook)
	auth.PUT("/books/:id", handlers.UpdateBook)
	auth.PATCH("/books/:id", handlers.PatchBook)
	auth.DELETE("/books/:id", handlers.DeleteBook)

	auth.POST("/categories", handlers.CreateCategory)
	auth.GET("/me", handlers.Profile)

	wishlist := r.Group("/wishlist")
	wishlist.Use(middlewares.JWTAuthMiddleware())

	wishlist.POST("/:bookId", handlers.AddToWishlist)
	wishlist.GET("/", handlers.GetWishlist)
	wishlist.DELETE("/:bookId", handlers.RemoveFromWishlist)

	cart := r.Group("/cart")
	cart.Use(middlewares.JWTAuthMiddleware())
	cart.POST("", handlers.AddToCart)
	cart.GET("", handlers.GetCart)
	cart.DELETE("/:id", handlers.RemoveFromCart)

	// run server
	port := os.Getenv("PORT") // Biar bisa ambil port dari Render
	if port == "" {
		port = "8080" // fallback local
	}
	r.Run(":" + port)
}
