package handlers

import (
	"net/http"
	"testing"

	"github.com/thylong/regexrace/models"

	"context"
)

func TestMgoDBFromR(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = req.WithContext(
		context.WithValue(req.Context(), "db", &models.FakeSession{}))

	fakeDb := MgoDBFromR(req)
	if fakeDb == nil {
		t.Error(req)
	}

}
