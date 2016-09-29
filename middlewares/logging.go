package middlewares

import (
	"net/http"
	"os"
	"time"

	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zenazn/goji/web/mutil"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

const logKey key = iota

// LoggingHandler set up Fields reused by every log instance.
func LoggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := viper.GetString("ENV")
		role := viper.GetString("ROLE")
		if env != "prod" {
			log.SetLevel(log.DebugLevel)
		}
		hostname, _ := os.Hostname()

		contextLogger := log.WithFields(log.Fields{
			"env":      env,
			"role":     role,
			"hostname": hostname,
		})

		if r != nil {
			r = r.WithContext(
				context.WithValue(r.Context(), logKey, *contextLogger))
		}
		next.ServeHTTP(w, r)
	})
}

// FromContext gets the logger out of the context.
func FromContext(ctx context.Context) log.Entry {
	logger, _ := ctx.Value(logKey).(log.Entry)
	return logger
}

// FromRequest gets the logger from the request's context.
func FromRequest(r *http.Request) log.Entry {
	return FromContext(r.Context())
}

// AccessLogHandler log any HTTP request made to the application.
func AccessLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		// Sniff the status and content size for logging
		lw := mutil.WrapWriter(w)

		next.ServeHTTP(lw, r)
		t2 := time.Now()

		logger := FromRequest(r)
		logger.WithFields(log.Fields{
			"method":      r.Method,
			"size":        lw.BytesWritten(),
			"uri":         r.RequestURI,
			"duration":    t2.Sub(t1),
			"status_code": lw.Status(),
		}).Info()
	})
}
