package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticHandler(t *testing.T) {
	cases := []struct {
		staticURLPath      string
		staticRelativePath string
	}{
		{"/static/robots.txt", "../static/robots.txt"},
		{"/static/regex.css", "../static/regex.css"},
		{"/static/regex.js", "../static/regex.js"},
	}

	for _, tc := range cases {
		req, err := http.NewRequest("GET", tc.staticURLPath, nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(StaticHandler)

		handler.ServeHTTP(w, req)

		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		robotsContent, err := ioutil.ReadFile(tc.staticRelativePath)
		if err != nil {
			t.Fatal(err)
		}

		expected := string(robotsContent)
		if w.Body.String() != expected {
			t.Errorf("handler returned unexpected body")
		}
	}
}
