package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/auth", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthHandler)

	handler.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var resMap map[string]string
	content, err := ioutil.ReadAll(w.Body)

	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Fatal(err)
	}
	err = json.Unmarshal(content, &resMap)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := resMap["token"]; !ok {
		t.Fatal(err)
	}
}
