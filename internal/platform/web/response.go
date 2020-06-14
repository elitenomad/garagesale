package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func Respond(w http.ResponseWriter, val interface{}, statusCode int) error  {
	data, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "Cannot marshall the data...")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		return  errors.Wrap(err, "Problem in processing the data to client...")
	}

	return nil
}
