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
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func LoginHandler(c *gin.Context) {
	var incomingUser models.UserRequest
	var dbUser models.User
	if err := c.ShouldBindJSON(&incomingUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	collection := models.Collection["users"]
	if err := collection.FindOne(c, bson.D{primitive.E{Key: "email", Value: incomingUser.Email}}).Decode(&dbUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(incomingUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := GenerateToken(dbUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error generating token: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func PatientRegisterHandler(c *gin.Context) {
	var user models.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	collection := models.Collection["users"]
	var existingUser models.UserRequest
	if err = collection.FindOne(c, bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&existingUser); err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "email already exists"})
		return
	}

	var newUser models.User
	if err := copier.Copy(&newUser, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("copier error: %v", err)})
		return
	}

	now := time.Now().Local()
	newUser.Password = hashedPassword
	newUser.Role = constant.PATIENT_ROLE
	newUser.CreatedAt = now
	newUser.UpdatedAt = now
	_, err = collection.InsertOne(c, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registration successful"})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateToken(userData models.User) (string, error) {
	expirationTime := time.Now().Local().Add(30 * time.Hour).Unix()
	claims := &Claims{
		UserID: userData.ID.Hex(),
		Email:  userData.Email,
		Role:   userData.Role,
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
