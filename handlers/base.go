package handlers

import (
	"errors"
	"net/http"

	"context"

	"github.com/spf13/viper"
	"github.com/thylong/regexrace/models"
)

// ErrJSONPayloadEmpty is returned when the JSON payload is empty.
var ErrJSONPayloadEmpty = errors.New("JSON payload is empty")

// ErrJSONPayloadInvalidBody is returned when the JSON payload is fucked up.
var ErrJSONPayloadInvalidBody = errors.New("Cannot parse request body")

// ErrJSONPayloadInvalidFormat is returned when the JSON payload is fucked up.
var ErrJSONPayloadInvalidFormat = errors.New("Invalid JSON format")

// MgoSessionFromCtx takes a context argument and return the related *mgo.session.
func MgoSessionFromCtx(ctx context.Context) models.Session {
	mgoSession, _ := ctx.Value("db").(models.Session)
	return mgoSession
}

// MgoDBFromR takes a request argument and return the extracted *mgo.session.
func MgoDBFromR(r *http.Request) models.DataLayer {
	return MgoSessionFromCtx(r.Context()).DB(viper.GetString("ROLE"))
}
