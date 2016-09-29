package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/thylong/regexrace/models"
)

// ScoreHandler stores scores from the request.
func ScoreHandler(w http.ResponseWriter, r *http.Request) {
	score := extractScoreFromRequest(r)

	// Get original question related to the given answer.
	questionsCol := models.MgoSessionFromR(r).DB("regexrace").C("scores")
	_, err := questionsCol.Upsert(bson.M{"username": score.Username}, score)
	if err != nil {
		panic(err)
	}

	responseData := make(map[string]interface{})
	responseData["status"] = "success"
	data, _ := json.Marshal(responseData)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// extractScoreFromRequest validates JSON Payload and store the score.
func extractScoreFromRequest(r *http.Request) models.Score {
	score := models.Score{}

	content, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(ErrJSONPayloadInvalidBody)
	}
	if len(content) == 0 {
		panic(ErrJSONPayloadEmpty)
	}

	err = json.Unmarshal(content, &score)
	if err != nil {
		panic(ErrJSONPayloadInvalidFormat)
	}
	return score
}
