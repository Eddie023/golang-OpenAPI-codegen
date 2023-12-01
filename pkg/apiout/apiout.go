// package apiout provides utility functions to handle API response.
package apiout

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func JSON(ctx context.Context, w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("marshalling json", err)
	}

	if _, err := w.Write(jsonData); err != nil {
		slog.Error("writing response", err)
	}
}

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func Error(ctx context.Context, w http.ResponseWriter, err error) {
	var er ErrorResponse
	var aerr *APIError
	switch {
	case IsBadRequest(err):

		w.WriteHeader(http.StatusBadRequest)
	case errors.As(err, &aerr):
		er = ErrorResponse{
			Error: err.Error(),
		}

		w.WriteHeader(aerr.Status)
	default:
		er = ErrorResponse{
			Error: http.StatusText(http.StatusInternalServerError),
		}
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(er)
}

// DecodeJSONBody decodes a JSON body and returns client-friendly errors.
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		err := errors.New("Content-Type header is not application/json")
		return NewRequestError(err, http.StatusUnsupportedMediaType)
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			err := fmt.Errorf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return NewRequestError(err, http.StatusBadRequest)

		case errors.Is(err, io.ErrUnexpectedEOF):
			err := fmt.Errorf("request body contains badly-formed JSON")
			return NewRequestError(err, http.StatusBadRequest)

		case errors.As(err, &unmarshalTypeError):
			err := fmt.Errorf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return NewRequestError(err, http.StatusBadRequest)

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			err := fmt.Errorf("request body contains unknown field %s", fieldName)
			return NewRequestError(err, http.StatusBadRequest)

		case errors.Is(err, io.EOF):
			err := errors.New("request body must not be empty")
			return NewRequestError(err, http.StatusBadRequest)

		case err.Error() == "http: request body too large":
			err := errors.New("request body must not be larger than 1MB")
			return NewRequestError(err, http.StatusRequestEntityTooLarge)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		err := errors.New("request body must only contain a single JSON object")
		return NewRequestError(err, http.StatusBadRequest)
	}

	return nil
}
