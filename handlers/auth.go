package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// AuthHandler returns a JWT token.
func AuthHandler(w http.ResponseWriter, r *http.Request) {
	signingKey := []byte(viper.GetString("TOKEN_SECRET"))

	expTime, _ := strconv.ParseInt(viper.GetString("TOKEN_TTL"), 10, 32)
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * time.Duration(expTime)).Unix(),
		Issuer:    viper.GetString("ROLE"),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, _ := token.SignedString(signingKey)

	responseData := make(map[string]interface{})
	responseData["token"] = signedString
	data, _ := json.Marshal(responseData)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
