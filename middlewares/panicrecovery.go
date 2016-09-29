package middlewares

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// PanicRecoveryHandler avoids application restarts in case of panic error.
func PanicRecoveryHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				status := 500

				if status >= 500 {
					log.Warn(err)
				}

				if err != nil {
					w.WriteHeader(500)

					if viper.GetString("ENV") == "dev" {
						fmt.Fprintf(w, "Error: "+err.(error).Error())
					}
					return
				}
				w.WriteHeader(status)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
