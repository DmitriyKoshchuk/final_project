package api

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("12345"))

func generateToken() (string, error) {
	// Используем хеш пароля в качестве claims
	hash := sha256.Sum256(jwtSecret)
	claims := jwt.MapClaims{
		"hash": hex.EncodeToString(hash[:]),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	// Проверяем хеш пароля
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		currentHash := sha256.Sum256(jwtSecret)
		return claims["hash"] == hex.EncodeToString(currentHash[:])
	}

	return false
}
