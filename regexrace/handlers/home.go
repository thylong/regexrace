package handlers

import "net/http"

// HomeHandler returns the Homepage and the first question.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/"+r.URL.Path[1:]+".html")
}
