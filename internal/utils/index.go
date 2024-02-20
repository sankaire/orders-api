package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type JWTClaims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type JWTPayload struct {
	ID    int    `json:"ID"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

var secretKey = []byte("secret-key")

func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}
func ComparePassword(hashedPassword string, textPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(textPassword))
	if err != nil {
		return false, err
	}
	return true, err
}
func CreateToken(payload JWTPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": payload.Email,
			"phone": payload.Phone,
			"id":    payload.ID,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return token, nil
}
