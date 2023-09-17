package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"patient/config"
	"patient/constant"
	"patient/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Claims struct to be encoded to JWT
type Claims struct {
	PatientId string `json:"patient_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	jwt.StandardClaims
}

func LoginHandler(c *gin.Context) {
	var incomingPatient models.PatientRequest
	var dbPatient models.Patient
	if err := c.ShouldBindJSON(&incomingPatient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	collection := models.Collection["patients"]
	if err := collection.FindOne(c, bson.D{primitive.E{Key: "email", Value: incomingPatient.Email}}).Decode(&dbPatient); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbPatient.Password), []byte(incomingPatient.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := GenerateToken(dbPatient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error generating token: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func PatientRegisterHandler(c *gin.Context) {
	var patient models.PatientRequest
	if err := c.ShouldBindJSON(&patient); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(patient.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	collection := models.Collection["patients"]
	var existingUser models.PatientRequest
	if err = collection.FindOne(c, bson.D{primitive.E{Key: "email", Value: patient.Email}}).Decode(&existingUser); err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email already exists"})
		return
	}

	var newPatient models.Patient
	if err := copier.Copy(&newPatient, &patient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("copier error: %v", err)})
		return
	}

	now := time.Now().Local()
	newPatient.Password = hashedPassword
	newPatient.Role = constant.PATIENT_ROLE
	newPatient.CreatedAt = now
	newPatient.UpdatedAt = now
	_, err = collection.InsertOne(c, newPatient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save patient: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateToken(userData models.Patient) (string, error) {
	expirationTime := time.Now().Local().Add(30 * time.Hour).Unix()
	claims := &Claims{
		PatientId: userData.ID.Hex(),
		Email:     userData.Email,
		Role:      userData.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Cfg.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRandomJWTKey() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("failed to generate random key: " + err.Error())
	}

	return base64.StdEncoding.EncodeToString(key)
}

// GetRequestBody reads the request body and returns it as a string.
// It also restores the request body back to its original state so it can be read again.
func GetRequestBody(c *gin.Context) (string, error) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			return "", err
		}
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes), nil
}
