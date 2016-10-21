package models

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	mgo "gopkg.in/mgo.v2"
)

type mongokey string

// MongoKey contains the Mongo session for the Request.
const MongoKey mongokey = "mongo"

// ErrNotFound returned when an object is not found.
var ErrNotFound = mgo.ErrNotFound

// MongoCollection wraps a mgo.Collection to embed methods in models.
type MongoCollection struct {
	*mgo.Collection
}

// Collection is an interface to access to the collection struct.
type Collection interface {
	Find(query interface{}) *mgo.Query
	Count() (n int, err error)
	FindId(id interface{}) *mgo.Query
	Insert(docs ...interface{}) error
	Remove(selector interface{}) error
	Update(selector interface{}, update interface{}) error
	Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error)
	EnsureIndex(index mgo.Index) error
	RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error)
}

// MongoDatabase wraps a mgo.Database to embed methods in models.
type MongoDatabase struct {
	*mgo.Database
}

// C shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (d MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

// DataLayer is an interface to access to the database struct
// (currently MongoDatabase).
type DataLayer interface {
	C(name string) Collection
	GetQuestions() ([]Question, error)
	GetQuestion(qid int) (Question, error)
	GetScores() ([]Score, error)
	FindTopScores() ([]Score, error)
}

// Session is an interface to access to the Session struct.
type Session interface {
	DB(name string) DataLayer
	SetSafe(safe *mgo.Safe)
	SetSyncTimeout(d time.Duration)
	SetSocketTimeout(d time.Duration)
	Close()
	Copy() Session
}

// MongoSession is currently a Mongo session.
type MongoSession struct {
	*mgo.Session
}

// DB shadows *mgo.DB to returns a DataLayer interface instead of *mgo.Database.
func (s MongoSession) DB(name string) DataLayer {
	return &MongoDatabase{Database: s.Session.DB(name)}
}

// Copy mocks mgo.Session.Copy()
func (s MongoSession) Copy() Session {
	return MongoSession{s.Session.Copy()}
}

// NewSession returns a new Mongo Session.
func NewSession() Session {
	mgoSession, err := mgo.Dial(viper.GetString("MONGO_URI"))
	session := MongoSession{mgoSession}
	if err != nil {
		panic(err)
	}
	session.SetSafe(&mgo.Safe{})
	session.SetSyncTimeout(3 * time.Second)
	session.SetSocketTimeout(3 * time.Second)

	return session
}

// PrepareDB ensure presence of persistent and immutable data in the DB.
func PrepareDB(session Session) {
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
