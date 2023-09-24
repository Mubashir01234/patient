package middleware

import (
	"errors"
	"net/http"
	"os"
	"patient/auth"
	"patient/constant"
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, authorizationTypeBearer) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		tokenStr := header[len(authorizationTypeBearer):]
		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token: " + err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
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
	if claims.Email == "" && claims.PatientId == "" && claims.Role == "" {
		return errors.New("unable to set claims")
	}

	c.Set(constant.PATIENT_ID_CONTEXT, claims.PatientId)
	c.Set(constant.EMAIL_CONTEXT, claims.Email)
	c.Set(constant.ROLE_CONTEXT, claims.Role)
	return nil
}
