package handlers

import (
	"errors"
	"net/http"

	"context"

	"github.com/spf13/viper"
	"github.com/thylong/regexrace/models"
	mgo "gopkg.in/mgo.v2"
)

// ErrJSONPayloadEmpty is returned when the JSON payload is empty.
var ErrJSONPayloadEmpty = errors.New("JSON payload is empty")

// ErrJSONPayloadInvalidBody is returned when the JSON payload is fucked up.
var ErrJSONPayloadInvalidBody = errors.New("Cannot parse request body")

// ErrJSONPayloadInvalidFormat is returned when the JSON payload is fucked up.
var ErrJSONPayloadInvalidFormat = errors.New("Invalid JSON format")

// MgoSessionFromCtx takes a context argument and return the related *mgo.session.
func MgoSessionFromCtx(ctx context.Context) *mgo.Session {
	mgoSession, _ := ctx.Value("db").(*mgo.Session)
	return mgoSession
}

// MgoSessionFromR takes a request argument and return the related *mgo.session.
func MgoSessionFromR(r *http.Request) *mgo.Session {
	return MgoSessionFromCtx(r.Context())
}

// MgoDBFromR takes a request argument and return the related *mgo.session.
func MgoDBFromR(r *http.Request) models.MongoDatabase {
	return models.MongoDatabase{Database: MgoSessionFromCtx(r.Context()).DB(viper.GetString("ROLE"))}
}
