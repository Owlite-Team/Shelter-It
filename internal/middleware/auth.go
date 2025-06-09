package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"strings"
	"time"
)

// AuthMiddleware verifies JWT token in incoming request
func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Extract token from Authorization header
		tokenString, err := extractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": getErrorMessage(err)})
			c.Abort()
			return
		}

		// Extract and validate claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Check token expiration
		if err := checkTokenExpiration(claims); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func extractTokenFromHeader(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Authorization format")
	}
	return parts[1], nil
}

func getErrorMessage(err error) string {
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		return "Invalid token signature"
	}
	return "Invalid or expired token"
}

func checkTokenExpiration(claims jwt.MapClaims) error {
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return errors.New("Token expired")
		}
	}
	return nil
}

func RateLimiter() gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second), 10)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
