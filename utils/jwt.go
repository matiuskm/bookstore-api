package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("UNKUHyBO3d1roAAdDHVLjnqBOpYb9SrI")

func GenerateJWT(username string, userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"userID":   userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(JwtKey)
}