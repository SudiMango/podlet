package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const maxJSONBodySize = 1 << 20

func SendJSON(w http.ResponseWriter, status int, response any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(response)
}

func SendError(w http.ResponseWriter, status int, err error) error {
	return SendJSON(w, status, map[string]string{
		"message": err.Error(),
	})
}

func DecodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, maxJSONBodySize))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}
	return nil
}
