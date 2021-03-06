package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func Respond(w http.ResponseWriter, val interface{}, statusCode int) error {
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	data, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "Cannot marshall the data...")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "Problem in processing the data to client...")
	}

	return nil
}

func RespondError(w http.ResponseWriter, err error) error {
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(w, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
