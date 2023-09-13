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
	Username string `json:"username"`
	jwt.StandardClaims
}

func LoginHandler(c *gin.Context) {
	var incomingUser models.UserRequest
	var dbUser models.User

	// Get JSON body
	if err := c.ShouldBindJSON(&incomingUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	collection := models.Collection["users"]
	// Fetch the user from the database
	if err := collection.FindOne(c, bson.D{primitive.E{Key: "email", Value: incomingUser.Email}}).Decode(&dbUser); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		}
		return
	}

	// if err := models.DB.Where("username = ?", incomingUser.Username).First(&dbUser).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	} else {
	// 	}
	// 	return
	// }

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(incomingUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}
	fmt.Println("------", dbUser)
	// Generate JWT token
	token, err := GenerateToken(dbUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error generating token: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func RegisterHandler(c *gin.Context) {
	var user models.UserRequest

	fmt.Println(GetRequestBody(c))

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	collection := models.Collection["users"]
	var existingUser models.UserRequest
	if err = collection.FindOne(c, bson.D{primitive.E{Key: "email", Value: user.Email}}).Decode(&existingUser); err == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email already exists"})
		return
	}
	var newUser models.User
	if err := copier.Copy(&newUser, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("copier error: %v", err)})
		return
	}
	now := time.Now().UTC()
	newUser.Password = hashedPassword
	newUser.CreatedAt = now
	newUser.UpdatedAt = now

	// Save the user to the database
	_, err = collection.InsertOne(c, newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Could not save user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateToken(email string) (string, error) {
	// The expiration time after which the token will be invalid.
	expirationTime := time.Now().Add(30 * time.Hour).Unix()

	// Create the JWT claims, which includes the email and expiration time
	claims := &jwt.StandardClaims{
		// In JWT, the expiry time is expressed as unix milliseconds
		ExpiresAt: expirationTime,
		Issuer:    email,
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(config.Cfg.JwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateRandomJWTKey() string {
	key := make([]byte, 32) // generate a 256 bit key
	_, err := rand.Read(key)
	if err != nil {
		panic("Failed to generate random key: " + err.Error())
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
	// Restore the io.ReadCloser to its original state
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// Convert to string and return
	return string(bodyBytes), nil
}
