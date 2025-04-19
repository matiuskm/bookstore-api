package handlers

import (
	"bookstore-api/database"
	"bookstore-api/models"
	"bookstore-api/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	var books []models.Book

	query := database.DB

	// get query params
	title := c.Query("title")
	author := c.Query("author")
	minPrice := c.Query("minPrice")
	maxPrice := c.Query("maxPrice")

	// filter by title
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	// filter by author
	if author != "" {
		query = query.Where("author LIKE ?", "%"+author+"%")
	}

	// filter by min price
	if minPrice != "" {
		query = query.Where("price >= ?", minPrice)
	}

	// filter by max price
	if maxPrice != "" {
		query = query.Where("price <= ?", maxPrice)
	}

	sort := c.Query("sort")
	order := c.Query("order")

	// sort by title, author, or price
	if sort != "" {
		if order == "desc" {
			order = "asc"
		}

		validSortFields := map[string]bool{
			"title":  true,
			"author": true,
			"price":  true,
		}

		if validSortFields[sort] {
			query = query.Order(fmt.Sprintf("%s %s", sort, order))
		}
	}

	if err := query.Preload("Category").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func GetBookByID(c *gin.Context) {
	id, _ := utils.GetUintID(c)
	var book []models.Book

	if err := database.DB.Preload("Category").Find(&book, id).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book data"})
		return
	}

	database.DB.Create(&book)
	c.JSON(http.StatusCreated, book)
}

func UpdateBook(c *gin.Context) {
	id,_ := utils.GetUintID(c)
	var book models.Book

	if err := database.DB.First(&book, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book data"})
		return
	}
	database.DB.Save(&book)
	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	id, _ := utils.GetUintID(c)

	database.DB.Delete(&models.Book{}, id)
	c.JSON(http.StatusNoContent, gin.H{"message": "Book deleted"})
}

func PatchBook(c *gin.Context) {
	id, err := utils.GetUintID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var book models.Book
	if err := database.DB.First(&book, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Gunakan map biar fleksibel dan bisa partial update
	var patchData map[string]interface{}
	if err := c.ShouldBindJSON(&patchData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cat, ok := patchData["categoryId"]; ok {
		// JSON numbers di-parse jadi float64
		catFloat, ok := cat.(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid categoryId"})
			return
		}
	
		// Convert ke uint
		catID := uint(catFloat)
	
		// Cek ke database
		var category models.Category
		if err := database.DB.First(&category, catID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
	
		// Valid → update field di patchData (optional, karena udah ada)
		patchData["category_id"] = catID
	}
	

	// Update pakai map
	if err := database.DB.Model(&book).Updates(patchData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Preload("Category").First(&book, book.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload book"})
		return
	}

	c.JSON(http.StatusOK, book)
}
