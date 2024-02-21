package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sankaire/orders-api/internal/db/internal/utils"
	"net/http"
	"strings"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		tokenString := request.Header.Get("Authorization")
		if tokenString == "" {
			utils.WriteResponse(response, http.StatusUnauthorized, false, "token not provided", nil)
			return
		}
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.WriteResponse(response, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}
		tokenString = parts[1]
		token, err := utils.VerifyToken(tokenString)
		if err != nil {
			utils.WriteResponse(response, http.StatusUnauthorized, false, "An error occurred", nil)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.WriteResponse(response, http.StatusUnauthorized, false, "Invalid token", nil)
			return
		}
		phone := claims["phone"]
		id := claims["id"]

		ctx := context.WithValue(request.Context(), "phone", phone)
		ctx = context.WithValue(ctx, "id", id)
		request = request.WithContext(ctx)

		next.ServeHTTP(response, request)
	})
}
