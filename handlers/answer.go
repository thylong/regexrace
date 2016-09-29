package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
	"github.com/thylong/regexrace/models"
)

// Answer .
type Answer struct {
	QID      int    `bson:"qid" json:"qid"`
	Regex    string `bson:"regex" json:"regex"`
	Modifier string `bson:"modifier" json:"modifier"`
}

// ErrJSONPayloadEmpty is returned when the JSON payload is empty.
var ErrJSONPayloadEmpty = errors.New("JSON payload is empty")

// ErrJSONPayloadInvalidBody is returned when the JSON payload is fucked up.
var ErrJSONPayloadInvalidBody = errors.New("Cannot parse request body")

// ErrJSONPayloadInvalidFormat is returned when the JSON payload is fucked up.
var ErrJSONPayloadInvalidFormat = errors.New("Invalid JSON format")

// AnswerHandler handler receive the JSON answer for a question_id and
// return JSON containing a status (fail|success) AND if success a new question.
func AnswerHandler(w http.ResponseWriter, r *http.Request) {
	answer := extractAnswerFromRequest(r)

	// Get original question related to the given answer.
	questionsCol := models.MgoSessionFromR(r).DB("regexrace").C("questions")
	var originalQuestion models.Question
	err := questionsCol.Find(bson.M{"qid": answer.QID}).One(&originalQuestion)
	if err != nil {
		panic(err)
	}

	responseData := make(map[string]interface{})
	if isAnswerRegexMatchQuestion(answer, originalQuestion) {
		responseData["status"] = "success"
		responseData["new_question"] = getNewJSONQuestion(answer, questionsCol)
	} else {
		responseData["status"] = "fail"
	}
	data, _ := json.Marshal(responseData)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// extractAnswerFromRequest validates JSON Payload and return answer.
func extractAnswerFromRequest(r *http.Request) Answer {
	answer := Answer{Modifier: "g"}

	content, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(ErrJSONPayloadInvalidBody)
	}
	if len(content) == 0 {
		panic(ErrJSONPayloadEmpty)
	}
	test := string(content)
	fmt.Println(test)
	err = json.Unmarshal(content, &answer)
	if err != nil {
		panic(ErrJSONPayloadInvalidFormat)
	}
	return answer
}

// getNewQuestion returns a new question with formatted HTML Sentence from the database.
func getNewQuestion(answer Answer, questionsCol *mgo.Collection) models.Question {
	var newQuestion models.Question

	err := questionsCol.Find(
		bson.M{"qid": answer.QID + 1},
	).Select(bson.M{"sentence": 1, "qid": 1, "match_positions": 1, "_id": 0}).One(&newQuestion)
	if err != nil {
		panic(err)
	}
	newQuestion.Sentence = models.FormatHTMLSentence(newQuestion.Sentence, newQuestion.MatchPositions)
	return newQuestion
}

// getNewJSONQuestion returns a new JSON question with formatted HTML Sentence from the database.
func getNewJSONQuestion(answer Answer, questionsCol *mgo.Collection) map[string]interface{} {
	newQuestion := getNewQuestion(answer, questionsCol)

	JSONQuestion := make(map[string]interface{})
	JSONQuestion["qid"] = newQuestion.QID
	JSONQuestion["sentence"] = newQuestion.Sentence
	JSONQuestion["match_positions"] = newQuestion.MatchPositions

	return JSONQuestion
}

// isAnswerRegexMatchQuestion returns true if the regex is a right answer else returns false.
func isAnswerRegexMatchQuestion(answer Answer, originalQuestion models.Question) bool {
	var re, err = regexp.Compile(answer.Regex)
	if err != nil {
		log.Warn(err.Error())
		return false
	}
	var matchPositions interface{}
	if answer.Modifier == "g" || answer.Modifier == "" {
		matchPositions = re.FindAllStringIndex(originalQuestion.Sentence, -1)
	} else {
		matchPositions = [][]int{re.FindStringIndex(originalQuestion.Sentence)}
	}

	fmt.Println(matchPositions)
	fmt.Println(originalQuestion.MatchPositions)
	if reflect.DeepEqual(matchPositions, originalQuestion.MatchPositions) {
		return true
	}
	return false
}
