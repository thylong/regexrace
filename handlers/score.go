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
	score := extractScoreFromRequest(r)
	token, err := middlewares.FromAuthHeader(r)
	err = score.SubmitScore(token)
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
	score := models.Score{Db: MgoDBFromR(r)}

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
