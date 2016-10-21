package models

import (
	"encoding/json"
	"io/ioutil"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// FakeDatabase satisfies DataLayer and act as a mock.
type FakeDatabase struct{}

// FakeCollection satisfies Collection and act as a mock.
type FakeCollection struct{}

// Find mock.
func (fc FakeCollection) Find(query interface{}) *mgo.Query {
	return nil
}

// FindId mock.
func (fc FakeCollection) FindId(id interface{}) *mgo.Query {
	return nil
}

// Count mock.
func (fc FakeCollection) Count() (n int, err error) {
	return 10, nil
}

// Insert mock.
func (fc FakeCollection) Insert(docs ...interface{}) error {
	return nil
}

// Remove mock.
func (fc FakeCollection) Remove(selector interface{}) error {
	return nil
}

// Update mock.
func (fc FakeCollection) Update(selector interface{}, update interface{}) error {
	return nil
}

// Upsert mock.
func (fc FakeCollection) Upsert(selector interface{}, update interface{}) (info *mgo.ChangeInfo, err error) {
	return nil, nil
}

// EnsureIndex mock.
func (fc FakeCollection) EnsureIndex(index mgo.Index) error {
	return nil
}

// RemoveAll mock.
func (fc FakeCollection) RemoveAll(selector interface{}) (info *mgo.ChangeInfo, err error) {
	return nil, nil
}

// C mocks mgo.Database(name).Collection(name).
func (db FakeDatabase) C(name string) Collection {
	return FakeCollection{}
}

// FakeSession satisfies Session and act as a mock of *mgo.session.
type FakeSession struct{}

// NewFakeSession mosck NewSession.
func NewFakeSession() Session {
	return FakeSession{}
}

// Close mocks mgo.Session.Close().
func (fs FakeSession) Close() {}

// Copy mocks mgo.Session.Copy().
// Regarding the context of use, no need to actually Copy the mock.
func (fs FakeSession) Copy() Session {
	return fs
}

// DB mocks mgo.Session.DB().
func (fs FakeSession) DB(name string) DataLayer {
	fakeDatase := FakeDatabase{}
	return fakeDatase
}

// SetSafe mocks mgo.Session.SetSafe().
func (fs FakeSession) SetSafe(safe *mgo.Safe) {}

// SetSyncTimeout mocks mgo.Session.SetSyncTimeout().
func (fs FakeSession) SetSyncTimeout(d time.Duration) {}

// SetSocketTimeout mocks mgo.Session.SetSocketTimeout().
func (fs FakeSession) SetSocketTimeout(d time.Duration) {}

// GetQuestions mocks models.GetQuestions().
func (db FakeDatabase) GetQuestions() ([]Question, error) {
	var Questions []Question
	questionContent, _ := ioutil.ReadFile(
		"/go/src/github.com/thylong/regexrace/config/default_questions.json")
	json.Unmarshal(questionContent, &Questions)

	return Questions, nil
}

// GetQuestion mocks models.GetQuestion().
func (db FakeDatabase) GetQuestion(qid int) (Question, error) {
	var Questions []Question
	questionContent, _ := ioutil.ReadFile(
		"/go/src/github.com/thylong/regexrace/config/default_questions.json")
	json.Unmarshal(questionContent, &Questions)

	for _, question := range Questions {
		if question.QID == qid {
			return question, nil
		}
	}
	return Question{}, nil
}

// GetScores mocks models.GetScores().
func (db FakeDatabase) GetScores() ([]Score, error) {
	var Scores []Score
	scoreContent, _ := ioutil.ReadFile(
		"/go/src/github.com/thylong/regexrace/config/default_scores.json")
	json.Unmarshal(scoreContent, &Scores)

	return Scores, nil
}

// FindTopScores mocks models.FindTopScores().
func (db FakeDatabase) FindTopScores() ([]Score, error) {
	var Scores []Score
	scoreContent, _ := ioutil.ReadFile(
		"/go/src/github.com/thylong/regexrace/config/default_scores.json")
	json.Unmarshal(scoreContent, &Scores)

	return Scores, nil
}
