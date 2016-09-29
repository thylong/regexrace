package models

import mgo "gopkg.in/mgo.v2"

// Score represent a unique score.
type Score struct {
	mgo.Collection
	// ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username  string `bson:"username" json:"username"`
	BestScore int    `bson:"best_score" json:"best_score"`
}
