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

	if len(patient.Email) > 0 && dbPatient.Email != patient.Email {
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
