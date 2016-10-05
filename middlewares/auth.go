package middlewares

import (
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// WithAuth filter requests depending on JWT token validity.
func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := FromAuthHeader(r)
		if err != nil || token == "" || !IsValidAuth(token) {
			w.WriteHeader(401)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// IsValidAuth returns true if the token submitted is valid otherwise false.
func IsValidAuth(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(viper.GetString("TOKEN_SIGNATURE")), nil
	})

	if err != nil {
		return false
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true
	}
	return false
}
