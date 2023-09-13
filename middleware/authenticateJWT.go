package middleware

import (
	"errors"
	"net/http"
	"patient/auth"
	"patient/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	tokenHeaderKey          = "Authorization"
	authorizationTypeBearer = "Bearer "
)

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(tokenHeaderKey)
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, authorizationTypeBearer) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		tokenStr := header[len(authorizationTypeBearer):]
		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return config.Cfg.JwtSecret, nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if err := setClaim(c, claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}

func setClaim(c *gin.Context, claims *auth.Claims) error {
	if claims.Username == "" {
		return errors.New("Unable to set username in claims")
	}

	c.Set("username", claims.Username)
	return nil
}
