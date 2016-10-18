package handlers

import (
	"net/http"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

// HomeHandler returns the Homepage and the first question.
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("home.html").ParseFiles(os.Getenv("GOPATH") + "/src/github.com/thylong/regexrace/static/home.html"))

	db := MgoDBFromR(r)
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
	t.Execute(w, data)
}
