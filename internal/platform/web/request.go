package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

func Decode(r *http.Request, val interface{}) error  {
	if err := json.NewDecoder(r.Body).Decode(val); err != nil {
		return errors.Wrap(err, "Decoding problem")
	}

	return nil
}
