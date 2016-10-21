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

func TestAnswerHandler(t *testing.T) {
	cases := []struct {
		HTTPVerb string
		Body     io.Reader
		Expected string
	}{
		{"POST", strings.NewReader(""), ErrJSONPayloadEmpty.Error()},
		{"POST", strings.NewReader("{wrong_format"), ErrJSONPayloadInvalidFormat.Error()},
		{"POST", bytes.NewBuffer([]byte(`{"q":1,"rex":"Hello","mr":""}`)), ErrJSONPayloadInvalidFormat.Error()},
		{"POST", bytes.NewBuffer([]byte(`{"qid":1,"regex":"Hello","modifier":""}`)), ""},
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

			req, err := http.NewRequest(tc.HTTPVerb, "/answer", tc.Body)
			req = req.WithContext(
				context.WithValue(req.Context(), "db", models.FakeSession{}),
			)

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(AnswerHandler)

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
