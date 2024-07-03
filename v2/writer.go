package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func getUnsupportedType(err error) string {
	marshalErr, ok := err.(*json.UnsupportedTypeError)
	if !ok {
		return err.Error()
	}

	return marshalErr.Type.String()
}

func indent(w http.ResponseWriter, status int, data any) error {
	payload, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		var unsupportedType *json.UnsupportedTypeError
		switch {
		case errors.As(err, &unsupportedType):
			return &Err{Code: http.StatusInternalServerError, message: fmt.Sprintf("payload contains unsupported type %s", getUnsupportedType(err))}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(payload)

	return nil
}

func (p *Parser) Write(w http.ResponseWriter, status int, data any) error {
	if status < 100 || status > 511 {
		status = 200
	}

	if p.Indent {
		return indent(w, status, data)
	}

	payload, err := json.Marshal(data)
	if err != nil {
		var unsupportedType *json.UnsupportedTypeError
		switch {
		case errors.As(err, &unsupportedType):
			return &Err{Code: http.StatusInternalServerError, message: fmt.Sprintf("payload contains unsupported type %s", getUnsupportedType(err))}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(payload)

	return nil
}
