package handlers

import (
	"net/http"
	"os"
)

// RobotsHandler servers static files.
func RobotsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, os.Getenv("GOPATH")+"/src/github.com/thylong/regexrace/static/robots.txt")
}
