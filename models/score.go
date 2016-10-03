package models

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Score represent a unique score.
type Score struct {
	MongoDatabase
	// ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username  string `bson:"username" json:"username"`
	BestScore int    `bson:"best_score" json:"best_score"`
}

// UpsertScore store/replace a score.
func (score *Score) UpsertScore() error {
	db := DB()
	_, err := db.C("scores").Upsert(bson.M{"username": score.Username}, score)
	if err != nil {
		log.Warning(err)
		return err
	}
	return nil
}

// EnsureScoreData is used to make sure the question & score collections are ready.
// The RemoveAll -> Insert is rough but will work at this point
// (TODO: Find a beautiful way to write this + Improve to do a smart insert)
func EnsureScoreData(session *mgo.Session) {
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

	scoreCol := session.DB("regexrace").C("scores")
	scoreCol.RemoveAll(bson.M{})
	scores := make([]interface{}, len(Scores))
	for i, v := range Scores {
		scores[i] = v
	}
	err = scoreCol.Insert(scores...)
	if err != nil {
		panic(err)
	}
	log.Info("Ensured Scores integrity.")
}
