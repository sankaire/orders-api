package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	url2 "net/url"
	"os"
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
type ResponseEnvelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var secretKey = []byte(os.Getenv("SECRET_KEY"))

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
func VerifyToken(tokenString string) (*jwt.Token, error) {
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

func WriteResponse(w http.ResponseWriter, statusCode int, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	envelope := ResponseEnvelope{
		Success: success,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(envelope)
}
func SendSms(phone string, item string, amount int64) error {
	url := "https://api.africastalking.com/version1/messaging"
	err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//
	//}
	apiKey := os.Getenv("AT_API_KEY")
	userName := os.Getenv("AT_USERNAME")

	params := url2.Values{
		"username": {userName},
		"to":       {phone},
		"message":  {fmt.Sprintf("Your order of %s has been placed successfully. You will pay KES %d on delivery.", item, amount)},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(params.Encode())))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	setHeaders(req, apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Unexpected response status:", resp.Status)
		return err
	}
	return nil
}
func setHeaders(request *http.Request, apiKey string) {
	request.Header.Set("apiKey", apiKey)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}
