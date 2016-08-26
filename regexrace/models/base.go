package models

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"context"

	"github.com/spf13/viper"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongokey string

// MongoKey contains the Mongo session for the Request.
const MongoKey mongokey = "mongo"

// ErrNotFound returned when an object is not found.
var ErrNotFound = mgo.ErrNotFound

// PrepareDB ensure presence of persistent and immutable data in the DB.
func PrepareDB(session *mgo.Session) {
	index := mgo.Index{
		Key:        []string{"qid", "sentence", "match_positions"},
		Unique:     true,
		DropDups:   true,
		Background: false,
	}
	err := session.DB("regexrace").C("regex").EnsureIndex(index)
	if err != nil {
		panic("Cannot ensure index ")
	}
}

// EnsureData is used to make sure the collection reflects questions.json content.
// The RemoveAll -> Insert is rough but will work at this point
// (TODO: Improve to do a smart insert)
func EnsureData(session *mgo.Session) {
	var Questions []Question

	fileContent, err := ioutil.ReadFile(
		"/go/src/github.com/thylong/regexrace/models/questions.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileContent, &Questions)
	if err != nil {
		panic(err)
	}

	collection := session.DB("regexrace").C("questions")
	Docsum, _ := collection.Count()

	if len(Questions) != Docsum {
		collection.RemoveAll(bson.M{})

		// This convert the []Regex slice to []interface slice
		// because of Insert() requirements.
		regexes := make([]interface{}, len(Questions))
		for i, v := range Questions {
			regexes[i] = v
		}
		err = collection.Insert(regexes...)
		if err != nil {
			panic(err)
		}
	}
}

// OpenMongoDB takes a mongo_uri argument and returns a mgosession or panics.
// TODO: As soon as a lazy mode is realeased on mgo, implement it.
func OpenMongoDB(mongoURI string) *mgo.Session {
	session, err := mgo.Dial(viper.GetString("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	session.SetSafe(&mgo.Safe{})
	session.SetSyncTimeout(3 * time.Second)
	session.SetSocketTimeout(3 * time.Second)

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
