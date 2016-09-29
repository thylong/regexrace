package middlewares

import (
	"net/http"
	"time"
)

// TimeoutHandler end the request after 2 seconds.
func TimeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, 3*time.Second, "timed out")
}
