package handlers

import "net/http"

// RobotsHandler servers static files.
func RobotsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./../static/robots.txt")
}
