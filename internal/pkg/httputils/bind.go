package httputils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func BindJSON(r *http.Request, dest interface{}) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("content-type must be application/json")
	}

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return errors.New("invalid JSON format")
	}

	return nil
}
