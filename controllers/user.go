package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"patient/auth"
	"patient/constant"
	"patient/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func GetPatientByEmail(c *gin.Context) {
	incomingEmail := c.Query("email")
	if len(incomingEmail) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	patientEmail, _ := c.Get(constant.EMAIL_CONTEXT)
	userRole, _ := c.Get(constant.ROLE_CONTEXT)
	if patientEmail != incomingEmail && userRole != constant.ADMIN_ROLE {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have no access to get user information"})
		return
	}

	var user models.GetPatient
	patientCollection := models.Collection["patients"]
	if err := patientCollection.FindOne(c, bson.D{primitive.E{Key: "email", Value: incomingEmail}}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, gin.H{"message": "no data present for this email", "data": nil})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UpdatePatient(c *gin.Context) {
	var patient models.UpdatePatientRequest
	var dbPatient models.Patient
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	patientEmail, _ := c.Get(constant.EMAIL_CONTEXT)
	collection := models.Collection["patients"]
	if err := collection.FindOne(c, bson.D{primitive.E{Key: "email", Value: patientEmail}}).Decode(&dbPatient); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "patient does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	if len(patient.Email) > 0 {
		err := collection.FindOne(c, bson.D{primitive.E{Key: "email", Value: patient.Email}})
		if err.Err() != nil && !errors.Is(err.Err(), mongo.ErrNoDocuments) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		dbPatient.Email = patient.Email
	}
	if len(patient.FullName) > 0 {
		dbPatient.FullName = patient.FullName
	}
	if len(patient.DOB) > 0 {
		dbPatient.DOB = patient.DOB
	}
	if len(patient.HomeAddress) > 0 {
		dbPatient.HomeAddress = patient.HomeAddress
	}
	if len(patient.MobileNumber) > 0 {
		dbPatient.MobileNumber = patient.MobileNumber
	}
	if len(patient.Password) > 0 {
		hashedPassword, err := auth.HashPassword(patient.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}
		dbPatient.Password = hashedPassword
	}
	dbPatient.UpdatedAt = time.Now().Local()

	res, err := collection.UpdateOne(c, bson.D{primitive.E{Key: "email", Value: patientEmail}}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "email", Value: dbPatient.Email},
				primitive.E{Key: "full_name", Value: dbPatient.FullName},
				primitive.E{Key: "dob", Value: dbPatient.DOB},
				primitive.E{Key: "home_address", Value: dbPatient.HomeAddress},
				primitive.E{Key: "phone", Value: dbPatient.MobileNumber},
				primitive.E{Key: "password", Value: dbPatient.Password},
				primitive.E{Key: "updated_at", Value: dbPatient.UpdatedAt},
			},
		},
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is already taken"})
		return
	}

	if res.MatchedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email does not exist"})
		return
	}

	token, err := auth.GenerateToken(dbPatient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error generating token: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "patient updated successfully",
		"token":   token,
	})
}

func DeleteBook(c *gin.Context) {
	incomingEmail := c.Param("email")
	if len(incomingEmail) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	patientEmail, _ := c.Get(constant.EMAIL_CONTEXT)
	userRole, _ := c.Get(constant.ROLE_CONTEXT)
	if patientEmail != incomingEmail && userRole != constant.ADMIN_ROLE {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have no access to get user information"})
		return
	}

	var user models.GetPatient
	patientCollection := models.Collection["patients"]
	if err := patientCollection.FindOne(c, bson.D{primitive.E{Key: "email", Value: incomingEmail}}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, gin.H{"message": "no data present for this email", "data": nil})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
		return
	}

	_, err := patientCollection.DeleteOne(c, bson.D{primitive.E{Key: "email", Value: incomingEmail}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on deleting patient: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "patient deleted successfully"})
}

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
