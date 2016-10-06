package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"context"

	mgo "gopkg.in/mgo.v2"
)

// TimeoutHandler end the request after 2 seconds.
func TimeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 3*time.Second, "timed out")
}

// FromAuthHeader is a "TokenExtractor" that takes a give request and extracts
// the JWT token from the Authorization header.
func FromAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

// DBFromContext gets the db out of the context.
func DBFromContext(ctx context.Context) *mgo.Session {
	db, _ := ctx.Value("db").(*mgo.Session)
	return db
}
