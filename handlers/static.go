package handlers

import (
	"net/http"
	"os"
)

// StaticHandler servers static files.
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	staticPath := os.Getenv("GOPATH") + "/src/github.com/thylong/regexrace" + r.URL.Path
	http.ServeFile(w, r, staticPath)
}
