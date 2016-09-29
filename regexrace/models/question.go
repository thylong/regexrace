package models

import mgo "gopkg.in/mgo.v2"

// Question represent a regex to find and the related context (sentence, match).
type Question struct {
	mgo.Collection
	QID            int     `bson:"qid" json:"qid"`
	Sentence       string  `bson:"sentence" json:"sentence"`
	MatchPositions [][]int `bson:"match_positions" json:"match_positions"`
	// Possibilities []interface{} // Valable answers (among others...)
}

// FormatHTMLSentence return a sentence with matches wrapped with HTML tags.
func FormatHTMLSentence(sentence string, MatchPositions [][]int) string {
	htmlSentence := sentence
	for index := len(MatchPositions) - 1; index >= 0; index-- {
		closingTagIndex := MatchPositions[index][1]
		openingTagIndex := MatchPositions[index][0]
		htmlSentence = htmlSentence[:openingTagIndex] + "<span class=\"highlighted\">" + htmlSentence[openingTagIndex:closingTagIndex] + "</span>" + htmlSentence[closingTagIndex:]
	}
	return htmlSentence
}
