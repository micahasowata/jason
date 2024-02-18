package jason

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

// Read parses the provided request body, triages possible errors and formats their message
func (j *Jason) Read(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	ok := isBodyJSON(r)
	if !ok {
		return &Err{Code: http.StatusUnsupportedMediaType, Msg: "content type is not application/json"}
	}

	ok = isDstPointer(dst)
	if !ok {
		return &Err{Code: http.StatusInternalServerError, Msg: "destination is not a pointer"}
	}

	r.Body = http.MaxBytesReader(w, r.Body, j.maxBodySize)

	dec := j.parser.NewDecoder(r.Body)

	if j.disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(&dst)
	if err != nil {
		switch {
		case isEmpty(err):
			return &Err{Code: http.StatusBadRequest, Msg: "request body must not be empty"}
		case containsImproperAssignment(err):
			return &Err{Code: http.StatusBadRequest, Msg: fmt.Sprintf("request body contains an invalid value for %s (at character %d)", getFieldName(getStmt(err.Error())), findErrorLocation(err.Error()))}
		case isBadlyFormedJSON(err):
			return &Err{Code: http.StatusBadRequest, Msg: fmt.Sprintf("request body contains badly formed JSON (at position %d)", findErrorLocation(err.Error()))}
		case j.disallowUnknownFields && isUnknownField(err):
			return &Err{Code: http.StatusBadRequest, Msg: fmt.Sprintf("request body contains unknown field %s", getFieldName(getStmt(err.Error())))}
		case isBodyTooLarge(err):
			return &Err{Code: http.StatusRequestEntityTooLarge, Msg: fmt.Sprintf("request body must not be larger than %d", j.maxBodySize)}
		default:
			return &Err{Code: http.StatusInternalServerError, Msg: err.Error()}
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return &Err{Code: http.StatusBadRequest, Msg: "request body must contain just one JSON object"}
	}

	return nil
}
