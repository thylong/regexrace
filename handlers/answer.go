package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/thylong/regexrace/models"
)

// Answer format.
type Answer struct {
	QID      int    `bson:"qid" json:"qid"`
	Regex    string `bson:"regex" json:"regex"`
	Modifier string `bson:"modifier" json:"modifier"`
}

// AnswerHandler handler receive the JSON answer for a question_id and
// return JSON containing a status (fail|success) AND if success a new question.
func AnswerHandler(w http.ResponseWriter, r *http.Request) {
	answer := extractAnswerFromRequest(r)
	originalQuestion := GetQuestion(answer.QID)

	responseData := make(map[string]interface{})
	if isAnswerMatchQuestion(answer, originalQuestion) {
		responseData["status"] = "success"
		responseData["new_question"] = originalQuestion.GetNextJSONQuestion(answer.QID)
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

	err = json.Unmarshal(content, &answer)
	if err != nil {
		panic(ErrJSONPayloadInvalidFormat)
	}
	return answer
}

// isAnswerMatchQuestion returns true if the regex is a right answer else returns false.
func isAnswerMatchQuestion(answer Answer, question models.Question) bool {
	var re = regexp.MustCompile(answer.Regex)
	submatches := make(map[int][][]int)

	matchPositions := re.FindAllStringSubmatchIndex(question.Sentence, -1)
	if answer.Modifier != "g" && answer.Modifier != "" {
		matchPositions = [][]int{matchPositions[0]}
	}
	submatches = splitFullMatchAndSubmatches(matchPositions, re.NumSubexp())

	log.Debug("MatchPositons retrieved: ", matchPositions)
	log.Debug("Submatches: ", submatches)
	log.Debug("MatchPositons expected: ", question.MatchPositions)
	for _, submatch := range submatches {
		if reflect.DeepEqual(submatch, question.MatchPositions) {
			return true
		}
	}

	return false
}

// splitFullMatchAndSubmatches return full match and each group in separated sub-arrays.
func splitFullMatchAndSubmatches(matchIndexes interface{}, numSub int) map[int][][]int {
	submatches := make(map[int][][]int)
	for _, subMatch := range matchIndexes.([][]int) {
		for num := 0; num <= numSub*2; num += 2 {
			extract := []int{subMatch[num], subMatch[num+1]}
			submatches[num] = append(submatches[num], extract)
		}
	}
	return submatches
}
