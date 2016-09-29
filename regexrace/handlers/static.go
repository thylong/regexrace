package handlers

import "net/http"

// StaticHandler servers static files.
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
