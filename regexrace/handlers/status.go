package handlers

import (
	"fmt"
	"net/http"
)

// StatusHandler endpoint to acknowledge application status.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
