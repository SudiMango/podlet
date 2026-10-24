package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/thecodearcher/limen"
)

const maxJSONBodySize = 1 << 20

func SendAPIResponse(
	responder *limen.Responder,
	w http.ResponseWriter,
	r *http.Request,
	response any,
	status ...int,
) error {
	if err, ok := response.(error); ok {
		return responder.Error(w, r, err)
	}

	responseStatus := http.StatusOK
	if len(status) > 0 {
		responseStatus = status[0]
	}

	return responder.JSON(w, r, responseStatus, response)
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
