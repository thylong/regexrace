package handlers

import (
	"net/http"
	"os"
	"text/template"
)

// LeaderboardHandler returns the Leaderboard with top 10 scores.
func LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(
		template.New("leaderboard.html").ParseFiles(os.Getenv("GOPATH") + "/src/github.com/thylong/regexrace/static/leaderboard.html"))

	db := MgoDBFromR(r)

	scores, err := db.FindTopScores()
	if err != nil {
		panic(err)
	}
	t.Execute(w, scores)
}
