package handlers

import (
	"bookstore-api/database"
	"bookstore-api/models"
	"bookstore-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddToWishlist(c *gin.Context) {
	userID := utils.GetUserID(c)
	bookIDParam := c.Param("bookId")
	bookID, _ := strconv.Atoi(bookIDParam)

	var book models.Book
	if err := database.DB.First(&book, bookID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Book not found"})
		return
	}

	var existing models.Wishlist
	err := database.DB.Where("user_id = ? AND book_id = ?", userID, bookID).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Book already in wishlist"})
		return
	}

	wishlist := models.Wishlist{
		UserID: userID,
		BookID: uint(bookID),
	}

	if err := database.DB.Create(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Book added to wishlist"})
}

func GetWishlist(c *gin.Context) {
	userID := utils.GetUserID(c)

	var wishlist []models.Wishlist
	if err := database.DB.Preload("Book.Category").Where("user_id = ?", userID).Find(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wishlist"})
		return
	}

	var books []models.Book
	for _, item := range wishlist {
		books = append(books, item.Book)
	}

	c.JSON(http.StatusOK, books)
}

func RemoveFromWishlist(c *gin.Context) {
	userID := utils.GetUserID(c)
	bookIDParam := c.Param("bookId")
	bookID, _ := strconv.Atoi(bookIDParam)

	var wishlist models.Wishlist
	if err := database.DB.Where("user_id = ? AND book_id = ?", userID, bookID).First(&wishlist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found in wishlist"})
		return
	}

	if err := database.DB.Delete(&wishlist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from wishlist"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Book removed from wishlist"})
}