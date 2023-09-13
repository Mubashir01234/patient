package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Healthcheck(g *gin.Context) {
	g.JSON(http.StatusOK, "ok")
}

// func FindBooks(c *gin.Context) {
// 	var books []models.Book

// 	// Get query params
// 	offsetQuery := c.DefaultQuery("offset", "0")
// 	limitQuery := c.DefaultQuery("limit", "10")

// 	// Convert query params to integers
// 	offset, _ := strconv.Atoi(offsetQuery)
// 	limit, _ := strconv.Atoi(limitQuery)

// 	// Create a cache key based on query params
// 	cacheKey := "books_offset_" + offsetQuery + "_limit_" + limitQuery

// 	// Try fetching the data from Redis first
// 	cachedBooks, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
// 	if err == nil {
// 		json.Unmarshal([]byte(cachedBooks), &books)
// 		c.JSON(http.StatusOK, gin.H{"data": books})
// 		return
// 	}

// 	// If cache missed, fetch data from the database
// 	models.DB.Offset(offset).Limit(limit).Find(&books)

// 	// Serialize books object and store it in Redis
// 	serializedBooks, _ := json.Marshal(books)
// 	cache.Rdb.Set(cache.Ctx, cacheKey, serializedBooks, 0)

// 	c.JSON(http.StatusOK, gin.H{"data": books})
// }

// func CreateBook(c *gin.Context) {
// 	var input models.CreateBook

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	book := models.Book{Title: input.Title, Author: input.Author}

// 	models.DB.Create(&book)

// 	c.JSON(http.StatusCreated, gin.H{"data": book})
// }

// func FindBook(c *gin.Context) {
// 	var book models.Book

// 	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": book})
// }

// func UpdateBook(c *gin.Context) {
// 	var book models.Book
// 	var input models.UpdateBook

// 	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	models.DB.Model(&book).Updates(models.Book{Title: input.Title, Author: input.Author})

// 	c.JSON(http.StatusOK, gin.H{"data": book})
// }

// func DeleteBook(c *gin.Context) {
// 	var book models.Book

// 	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
// 		return
// 	}

// 	models.DB.Delete(&book)

// 	c.JSON(http.StatusNoContent, gin.H{"data": true})
// }
