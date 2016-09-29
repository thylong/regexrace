package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"context"

	"github.com/spf13/viper"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongokey string

// MongoKey contains the Mongo session for the Request.
const MongoKey mongokey = "mongo"

// ErrNotFound returned when an object is not found.
var ErrNotFound = mgo.ErrNotFound

// PrepareDB ensure presence of persistent and immutable data in the DB.
func PrepareDB(session *mgo.Session) {
	indexes := []mgo.Index{
		mgo.Index{
			Key:        []string{"qid", "sentence", "match_positions"},
			Unique:     true,
			DropDups:   true,
			Background: false,
		},
		mgo.Index{
			Key:        []string{"username"},
			Unique:     true,
			DropDups:   true,
			Background: false,
		},
	}
	for _, index := range indexes {
		err := session.DB("regexrace").C("regex").EnsureIndex(index)
		if err != nil {
			panic("Cannot ensure index ")
		}
	}
	fmt.Println("Prepared database indexes.")
}

// EnsureData is used to make sure the question & score collections are ready.
// The RemoveAll -> Insert is rough but will work at this point
// (TODO: Find a beautiful way to write this + Improve to do a smart insert)
func EnsureData(session *mgo.Session) {
	// var collections map[string]mgo.Collection
	var Questions []Question
	var Scores []Score

	questionContent, err := ioutil.ReadFile(
		"/go/src/github.com/thylong/regexrace/models/questions.json")
	if err != nil {
		panic(err)
	}
	scoreContent, err := ioutil.ReadFile(
		"/go/src/github.com/thylong/regexrace/models/default_scores.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(questionContent, &Questions)
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
	fmt.Println("Ensured data integrity.")
}

// OpenMongoDB takes a mongo_uri argument and returns a mgosession or panics.
// See : http://stackoverflow.com/questions/26574594/best-practice-to-maintain-a-mgo-session
func OpenMongoDB() *mgo.Session {
	session := viper.Get("MONGO_SESSION").(*mgo.Session).Copy()

	return session
}

// MgoSessionFromCtx takes a context argument and return the related *mgo.session.
func MgoSessionFromCtx(ctx context.Context) *mgo.Session {
	mgoSession, _ := ctx.Value(MongoKey).(*mgo.Session)
	return mgoSession
}

// MgoSessionFromR takes a request argument and return the related *mgo.session.
func MgoSessionFromR(r *http.Request) *mgo.Session {
	return MgoSessionFromCtx(r.Context())
}
