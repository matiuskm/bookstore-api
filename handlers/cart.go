package handlers

import (
	"bookstore-api/database"
	"bookstore-api/models"
	"bookstore-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCart(c *gin.Context) {
	userID := utils.GetUserID(c)
	var cartItems []models.CartItem
	if err := database.DB.Preload("Book").Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cartItems})
}

func AddToCart(c *gin.Context) {
	userID := utils.GetUserID(c)

	var input struct {
		BookID 		uint `json:"bookId" binding:"required"`
		Quantity 	int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Quantity <= 0 {
		input.Quantity = 1
	}

	var cartItem models.CartItem
	err := database.DB.Where("user_id =? AND book_id =?", userID, input.BookID).First(&cartItem).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart item"})
		return
	}

	if err == gorm.ErrRecordNotFound {
		newCartItem := models.CartItem{
			UserID: userID,
			BookID: input.BookID,
			Quantity: input.Quantity,
		}

		if err := database.DB.Create(&newCartItem).Error; err!= nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"cartItem": newCartItem})
		return
	}

	cartItem.Quantity += input.Quantity
	if err := database.DB.Save(&cartItem).Error; err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cartItem": cartItem})
}

func RemoveFromCart(c *gin.Context) {
	userID := utils.GetUserID(c)
	cartItemID := c.Param("id")

	var cartItem models.CartItem
	if err := database.DB.Where("user_id =? AND id =?", userID, cartItemID).First(&cartItem).Error; err!= nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	if err := database.DB.Delete(&cartItem).Error; err!= nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}