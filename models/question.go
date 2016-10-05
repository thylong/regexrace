package models

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Question represent a regex to find and the related context (sentence, match).
type Question struct {
	Db             MongoDatabase
	QID            int     `bson:"qid" json:"qid"`
	Sentence       string  `bson:"sentence" json:"sentence"`
	MatchPositions [][]int `bson:"match_positions" json:"match_positions"`
	// Possibilities []interface{} // Valable answers (among others...)
}

// FormatHTMLSentence return a sentence with matches wrapped with HTML tags.
func (q *Question) FormatHTMLSentence() string {
	htmlSentence := q.Sentence

	for index := len(q.MatchPositions) - 1; index >= 0; index-- {
		closingTagIndex := q.MatchPositions[index][1]
		openingTagIndex := q.MatchPositions[index][0]
		htmlSentence = htmlSentence[:openingTagIndex] + "<span class=\"highlighted\">" + htmlSentence[openingTagIndex:closingTagIndex] + "</span>" + htmlSentence[closingTagIndex:]
	}
	return htmlSentence
}

// GetQuestions returns all Questions.
func (db *MongoDatabase) GetQuestions() ([]Question, error) {
	var questions []Question

	err := db.C("questions").Find(bson.M{}).All(&questions)
	if err != nil {
		return questions, err
	}
	return questions, nil
}

// GetQuestion returns a Question from an ID otherwise nil.
func (db *MongoDatabase) GetQuestion(qid int) (Question, error) {
	var originalQuestion Question

	err := db.C("questions").Find(bson.M{"qid": qid}).One(&originalQuestion)
	originalQuestion.Db = *db
	if err != nil {
		return originalQuestion, err
	}
	return originalQuestion, nil
}

// GetNextJSONQuestion returns a new JSON question with formatted HTML Sentence from the database.
func (q *Question) GetNextJSONQuestion(qid int) map[string]interface{} {
	newQuestion, _ := q.Db.GetQuestion(qid + 1)

	JSONQuestion := make(map[string]interface{})
	JSONQuestion["qid"] = newQuestion.QID
	JSONQuestion["sentence"] = newQuestion.FormatHTMLSentence()
	JSONQuestion["match_positions"] = newQuestion.MatchPositions

	return JSONQuestion
}

// EnsureQuestionData is used to make sure the question collection ready.
// The RemoveAll -> Insert is rough but will work at this point
// (TODO: Find a beautiful way to write this + Improve to do a smart insert)
func EnsureQuestionData(session *mgo.Session) {
	var Questions []Question

	questionContent, err := ioutil.ReadFile(
		"/go/src/github.com/thylong/regexrace/config/default_questions.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(questionContent, &Questions)
	if err != nil {
		panic(err)
	}

	questionCol := session.DB("regexrace").C("questions")
	Docsum, _ := questionCol.Count()
	if len(Questions) != Docsum {
		questionCol.RemoveAll(bson.M{})

		// This convert the []Regex slice to []interface slice
		// because of Insert() requirements.
		regexes := make([]interface{}, len(Questions))
		for i, v := range Questions {
			regexes[i] = v
		}

		err = questionCol.Insert(regexes...)
		if err != nil {
			panic(err)
		}
	}
	log.Info("Ensured Questions integrity.")
}
