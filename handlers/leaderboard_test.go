package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"context"

	"github.com/thylong/regexrace/models"
)

func TestLeaderboardHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/home", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(
		context.WithValue(req.Context(), "db", models.FakeSession{}))

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(LeaderboardHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	content, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Fatal(err)
	}
}
