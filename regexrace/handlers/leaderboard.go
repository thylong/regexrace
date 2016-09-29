package handlers

import (
	"net/http"
	"text/template"

	"gopkg.in/mgo.v2/bson"

	"github.com/thylong/regexrace/models"
)

// LeaderboardHandler returns the Leaderboard with top 10 scores.
func LeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// http.ServeFile(w, r, "static/"+r.URL.Path[1:]+".html")
	t := template.Must(template.New("leaderboard.html").ParseFiles("static/leaderboard.html"))

	var scores []models.Score

	scoresCol := models.MgoSessionFromR(r).DB("regexrace").C("scores")
	err := scoresCol.Find(bson.M{}).Sort("-best_score").All(&scores)
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, scores)
	if err != nil {
		panic(err)
	}
}
