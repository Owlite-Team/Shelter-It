package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func GenerateToken(id uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	// Check for error
	if err != nil {
		return nil, err
	}

	// Validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
