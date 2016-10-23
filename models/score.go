package models

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Score represent a unique score.
type Score struct {
	Username  string `bson:"username" json:"username"`
	BestScore int    `bson:"best_score" json:"best_score"`
	Submitted bool   `bson:"submitted" json:"submitted"`
}

// UpsertScore store/replace a score.
func (score *Score) UpsertScore(db DataLayer) error {
	_, err := db.C("scores").Upsert(
		bson.M{"username": score.Username}, score)
	if err != nil {
		log.Warning(err)
		return err
	}
	return nil
}

// GetScores returns all scores.
func (db *MongoDatabase) GetScores() ([]Score, error) {
	var scores []Score
	err := db.C("scores").Find(bson.M{}).All(&scores)
	if err != nil {
		log.Warning(err)
		return nil, err
	}
	return scores, nil
}

// FindTopScores returns all scores.
func (db *MongoDatabase) FindTopScores() ([]Score, error) {
	var scores []Score
	err := db.C("scores").Find(bson.M{"submitted": true}).Sort("-best_score").All(&scores)
	if err != nil {
		log.Warning(err)
		return nil, err
	}
	return scores, nil
}

// SubmitScore replace token by username and set submitted to true.
func (score *Score) SubmitScore(db DataLayer, token string) error {
	update := bson.M{"$set": bson.M{
		"username": score.Username, "submitted": true}}

	err := db.C("scores").Update(
		bson.M{"username": token}, update)
	if err != nil {
		log.Warning(err)
		return err
	}
	return nil
}

// InsertScore store a new score.
func (score *Score) InsertScore(db DataLayer) error {
	err := db.C("scores").Insert(score)
	if err != nil {
		log.Warning(err)
		return err
	}
	return nil
}

// EnsureScoreData makes sure the score collection is ready.
// The RemoveAll -> Insert is rough but will work at this point
// (TODO: Find a beautiful way to write this + Improve to do a smart insert)
func EnsureScoreData(session Session) {
	scoreCol := session.DB("regexrace").C("scores")
	Docsum, _ := scoreCol.Count()
	if Docsum <= 3 {
		var Scores []Score

		scoreContent, err := ioutil.ReadFile(
			"/go/src/github.com/thylong/regexrace/config/default_scores.json")
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(scoreContent, &Scores)
		if err != nil {
			panic(err)
		}
		scoreCol.RemoveAll(bson.M{})
		scores := make([]interface{}, len(Scores))
		for i, v := range Scores {
			scores[i] = v
		}
		err = scoreCol.Insert(scores...)
		if err != nil {
			panic(err)
		}
	}
	log.Info("Ensured Scores integrity.")
}
