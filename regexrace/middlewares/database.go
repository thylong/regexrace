package middlewares

import (
	"net/http"

	"github.com/thylong/regexrace/models"

	"context"
)

// MongoHandler insert Mgo.session in context and serve the request.
func MongoHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r != nil {
			mgoSession := models.OpenMongoDB()
			defer mgoSession.Close()
			r = r.WithContext(
				context.WithValue(r.Context(), models.MongoKey, mgoSession))
		}
		next.ServeHTTP(w, r)
	})
}
