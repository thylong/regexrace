package models

// Question represent a regex to find and the related context (sentence, match).
type Question struct {
	QID            int     `bson:"qid" json:"qid"`
	Sentence       string  `bson:"sentence" json:"sentence"`
	MatchPositions [][]int `bson:"match_positions" json:"match_positions"`
	// Possibilities []interface{} // Valable answers (among others...)
}
