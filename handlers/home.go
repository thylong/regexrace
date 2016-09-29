package handlers

import (
	"net/http"
	"text/template"

	"gopkg.in/mgo.v2/bson"

	"github.com/spf13/viper"
	"github.com/thylong/regexrace/models"
)

// HomeHandler returns the Homepage and the first question.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("home.html").ParseFiles("static/home.html"))

	var firstQuestion models.Question
	err := models.MgoSessionFromR(r).DB("regexrace").C("questions").Find(
		bson.M{"qid": 0}).One(&firstQuestion)
	if err != nil {
		panic(err)
	}
	htmlSentence := models.FormatHTMLSentence(firstQuestion.Sentence, firstQuestion.MatchPositions)

	data := struct {
		Sentence      string
		QID           int
		TimerDuration string
	}{
		Sentence:      htmlSentence,
		QID:           firstQuestion.QID,
		TimerDuration: viper.GetString("TIMER_DURATION"),
	}
	err = t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
