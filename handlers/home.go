package handlers

import (
	"net/http"
	"text/template"

	"github.com/spf13/viper"
	"github.com/thylong/regexrace/models"
)

// HomeHandler returns the Homepage and the first question.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("home.html").ParseFiles("static/home.html"))

	db := models.DB()
	firstQuestion, _ := db.GetQuestion(1)
	htmlSentence := firstQuestion.FormatHTMLSentence()

	data := struct {
		Sentence      string
		QID           int
		TimerDuration string
	}{
		Sentence:      htmlSentence,
		QID:           firstQuestion.QID,
		TimerDuration: viper.GetString("TIMER_DURATION"),
	}
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
