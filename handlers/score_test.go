package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"context"

	"github.com/thylong/regexrace/models"
)

func TestScoreHandler(t *testing.T) {
	cases := []struct {
		HTTPVerb string
		Body     io.Reader
		Expected string
	}{
		{"POST", strings.NewReader(""), ErrJSONPayloadEmpty.Error()},
		{"POST", strings.NewReader("{wrong_format"), ErrJSONPayloadInvalidFormat.Error()},
		{"POST", bytes.NewBuffer([]byte(`{"best":3,"user":"test", "toke":"test"}`)), ErrJSONPayloadInvalidFormat.Error()},
		{"POST", bytes.NewBuffer([]byte(`{"best_score":3,"username":"test", "token":"test"}`)), ""},
	}

	for _, tc := range cases {
		// This closure encapsulates the test to prevent panic to break the loop.
		func() {
			// Recovers from expected panics.
			defer func() {
				if err := recover(); err != nil {
					if tc.Expected != err.(error).Error() {
						t.Errorf("handler panic with wrong error : got %v want %v",
							err.(error).Error(), tc.Expected)
					}
				}
				return
			}()

			req, err := http.NewRequest(tc.HTTPVerb, "/score", tc.Body)
			req = req.WithContext(
				context.WithValue(req.Context(), "db", models.FakeSession{}),
			)

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(ScoreHandler)

			handler.ServeHTTP(w, req)

			if status := w.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}

			content, err := ioutil.ReadAll(w.Body)
			if err != nil || len(content) == 0 {
				t.Fatal(err)
			}
		}()
	}
}
