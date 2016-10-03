package handlers

import (
	"errors"
	"net/http"

	"context"

	"github.com/thylong/regexrace/models"
	mgo "gopkg.in/mgo.v2"
)

type mongokey string

// MongoKey contains the Mongo session for the Request.
const MongoKey mongokey = "mongo"

// ErrJSONPayloadEmpty is returned when the JSON payload is empty.
var ErrJSONPayloadEmpty = errors.New("JSON payload is empty")

// ErrJSONPayloadInvalidBody is returned when the JSON payload is fucked up.
var ErrJSONPayloadInvalidBody = errors.New("Cannot parse request body")

// ErrJSONPayloadInvalidFormat is returned when the JSON payload is fucked up.
var ErrJSONPayloadInvalidFormat = errors.New("Invalid JSON format")

// GetQuestion helps to get a question struct from a request context.
func GetQuestion(qid int) models.Question {
	db := models.DB()
	question, err := db.GetQuestion(qid)
	if err != nil {
		panic(err)
	}
	return question
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
