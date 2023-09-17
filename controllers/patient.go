package controllers

import (
	"encoding/json"
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
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetAllPatients(c *gin.Context) {
	userRole, _ := c.Get(constant.ROLE_CONTEXT)
	if userRole != constant.ADMIN_ROLE {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you have no access to get patient information"})
		return
	}

	filterJson := c.Query("json")
	filter := models.Filter{
		Page:  1,
		Limit: 10,
	}
	if len(filterJson) > 0 {
		if err := json.Unmarshal([]byte(filterJson), &filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid filter request passed"})
			return
		}
	}

	if filter.Page < 1 || filter.Limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filter values must be greater than 0"})
		return
	}

	countOpts := options.Count()
	collection := models.Collection["patients"]
	totalRecords, err := collection.CountDocuments(c, bson.M{}, countOpts)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to get record count: " + err.Error()})
		return
	}

	skip := (filter.Page - 1) * filter.Limit
	opts := options.Find()
	opts.SetSort(bson.D{primitive.E{Key: "updated_at", Value: -1}})
	opts.SetLimit(filter.Limit)
	opts.SetSkip(int64(skip))
	var result []models.GetPatient
	cursor, err := collection.Find(c, bson.M{}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, gin.H{"message": "no patient found", "data": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
		return
	}

	for cursor.Next(c) {
		var content models.GetPatient
		err := cursor.Decode(&content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error on getting data: " + err.Error()})
			return
		}

		result = append(result, content)
	}
	metadata := computeMetadata(totalRecords, filter.Page, filter.Limit)
	c.JSON(http.StatusOK, gin.H{
		"data":     result,
		"metadata": metadata,
	})
}
