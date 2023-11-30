package apiout

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeJSONBody(t *testing.T) {
	type testcase struct {
		name            string
		giveBody        string
		giveContentType string
		wantErr         error
	}

	testcases := []testcase{
		{name: "ok", giveBody: `{"test": "ok"}`, giveContentType: "application/json", wantErr: nil},
		{name: "no close bracket", giveBody: `{`, giveContentType: "application/json", wantErr: &APIError{Err: errors.New("request body contains badly-formed JSON"), Status: http.StatusBadRequest}},
		{name: "multiple objects", giveBody: `{"test": "ok"}{"second": "ok"}`, giveContentType: "application/json", wantErr: &APIError{Err: errors.New("request body must only contain a single JSON object"), Status: http.StatusBadRequest}},
		{name: "empty", giveBody: ``, giveContentType: "application/json", wantErr: &APIError{Err: errors.New("request body must not be empty"), Status: http.StatusBadRequest}},
		{name: "invalid content type", giveBody: `{"test": "ok"}`, giveContentType: "other", wantErr: &APIError{Err: errors.New("Content-Type header is not application/json"), Status: http.StatusUnsupportedMediaType}},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var i json.RawMessage
			w := httptest.NewRecorder()

			body := io.NopCloser(strings.NewReader(tc.giveBody))

			r := http.Request{
				Body:   body,
				Header: make(http.Header),
			}
			r.Header.Add("Content-Type", tc.giveContentType)

			err := DecodeJSONBody(w, &r, &i)

			if tc.wantErr == nil && err != nil {
				t.Fatalf("wanted no err but got %s", err)
			}

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}
		})
	}
}
