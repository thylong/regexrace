package models

import (
	"github.com/spf13/viper"

	log "github.com/Sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
)

type mongokey string

// MongoKey contains the Mongo session for the Request.
const MongoKey mongokey = "mongo"

// ErrNotFound returned when an object is not found.
var ErrNotFound = mgo.ErrNotFound

// MongoDatabase wraps a mgo.Database in order to embed methods in models.
type MongoDatabase struct {
	*mgo.Database
}

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
			panic("Cannot ensure index.")
		}
	}
	log.Info("Prepared database indexes.")
}

// DB takes a mongo_uri argument and returns a mgosession or panics.
// See : http://stackoverflow.com/questions/26574594/best-practice-to-maintain-a-mgo-session
func DB() MongoDatabase {
	session := viper.Get("MONGO_SESSION").(*mgo.Session).Copy()

	return MongoDatabase{Database: session.DB(viper.GetString("ROLE"))}
}
