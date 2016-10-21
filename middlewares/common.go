package middlewares

import (
	"net/http"
	"time"

	"context"

	mgo "gopkg.in/mgo.v2"
)

// TimeoutHandler end the request after 2 seconds.
func TimeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 3*time.Second, "timed out")
}

// DBFromContext gets the db out of the context.
func DBFromContext(ctx context.Context) *mgo.Session {
	db, _ := ctx.Value("db").(*mgo.Session)
	return db
}
