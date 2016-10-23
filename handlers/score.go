package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/thylong/regexrace/middlewares"
	"github.com/thylong/regexrace/models"
)

// ScoreHandler stores scores from the request.
func ScoreHandler(w http.ResponseWriter, r *http.Request) {
	score, err := extractScoreFromRequest(r)
	if err != nil {
		panic(err)
	}
	token, err := middlewares.FromAuthHeader(r)
	if err != nil {
		panic(err)
	}

	db := MgoDBFromR(r)
	err = score.SubmitScore(db, token)
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
func extractScoreFromRequest(r *http.Request) (models.Score, error) {
	var score = models.Score{}

	content, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if len(content) == 0 {
		return score, ErrJSONPayloadEmpty
	}

	err = json.Unmarshal(content, &score)
	if (err != nil || score == models.Score{Username: "", BestScore: 0}) {
		return score, ErrJSONPayloadInvalidFormat
	}
	return score, nil
}
