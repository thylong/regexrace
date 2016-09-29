package main

import (
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/spf13/viper"
	"github.com/thylong/regexrace/config"
	"github.com/thylong/regexrace/handlers"
	"github.com/thylong/regexrace/middlewares"
	"github.com/thylong/regexrace/models"

	"github.com/justinas/alice"
)

func main() {
	config.LoadConfig()

	// Create and store a Mongo session for every requests.
	session, err := mgo.Dial(viper.GetString("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	session.SetSafe(&mgo.Safe{})
	session.SetSyncTimeout(3 * time.Second)
	session.SetSocketTimeout(3 * time.Second)
	viper.Set("MONGO_SESSION", session)
	defer session.Close()

	// Ensure data looks like expected.
	models.PrepareDB(session)
	models.EnsureData(session)

	// Middlewares triggered for every requests.
	c := alice.New(
		middlewares.LoggingHandler,
		middlewares.TimeoutHandler,
		middlewares.AccessLogHandler,
		middlewares.MongoHandler,
	)
	if viper.GetString("ENV") != "dev" {
		c.Append(middlewares.PanicRecoveryHandler) // Has to be the latest middleware.
	}

	// Register Handlers.
	http.Handle("/status", c.ThenFunc(http.HandlerFunc(handlers.StatusHandler)))
	http.Handle("/home", c.ThenFunc(http.HandlerFunc(handlers.HomeHandler)))
	http.Handle("/leaderboard", c.ThenFunc(http.HandlerFunc(handlers.LeaderboardHandler)))
	http.Handle("/answer", c.ThenFunc(http.HandlerFunc(handlers.AnswerHandler)))
	http.Handle("/score", c.ThenFunc(http.HandlerFunc(handlers.ScoreHandler)))
	// Serve css and js
	http.Handle("/static/", c.ThenFunc(http.HandlerFunc(handlers.StaticHandler)))
	http.ListenAndServe(":8080", nil)
}
