package models

import (
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
	indexes := make(map[string]mgo.Index)
	indexes["questions"] = mgo.Index{
		Key:        []string{"qid", "sentence", "match_positions"},
		Unique:     true,
		DropDups:   true,
		Background: false,
	}
	indexes["scores"] = mgo.Index{
		Key:        []string{"username"},
		Unique:     false,
		DropDups:   false,
		Background: false,
	}
	for collectionName, index := range indexes {
		err := session.DB("regexrace").C(collectionName).EnsureIndex(index)
		if err != nil {
			panic("Cannot ensure index.")
		}
	}
	log.Info("Prepared database indexes.")
}
