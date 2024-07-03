package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

func isDstPointer(dst any) bool {
	return reflect.TypeOf(dst).Kind() == reflect.Ptr
}

func isBodyJSON(r *http.Request) bool {
	header := r.Header.Get("Content-Type")
	values := strings.Split(header, ";")
	mediaType := strings.ToLower(strings.TrimSpace(values[0]))
	return mediaType == ContentType
}

func getUnknownField(err error) string {
	errormessageSlice := strings.Split(err.Error(), " ")
	return errormessageSlice[len(errormessageSlice)-1]
}

func getTypeError(err error) string {
	unmarshalErr, ok := err.(*json.UnmarshalTypeError)
	if !ok {
		return err.Error()
	}

	if len(unmarshalErr.Field) > 0 {
		return fmt.Sprintf(`body contains invalid value type for "%s" at character %d`, unmarshalErr.Field, unmarshalErr.Offset)
	}

	return fmt.Sprintf("body contains invalid value type at character %d", unmarshalErr.Offset)
}

// Read reads values from input stream
func (p *Parser) Read(w http.ResponseWriter, r *http.Request, dst any) error {
	hasCorrectDst := isDstPointer(dst)
	if !hasCorrectDst {
		return &Err{Code: http.StatusInternalServerError, message: "destination is not a pointer"}
	}

	isJSON := isBodyJSON(r)
	if !isJSON {
		return &Err{Code: http.StatusUnsupportedMediaType, message: "request body content type is not application/json"}
	}

	r.Body = http.MaxBytesReader(w, r.Body, int64(p.Limit))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalError *json.UnmarshalTypeError

		switch {
		case errors.Is(err, io.EOF):
			return &Err{Code: http.StatusBadRequest, message: "body is empty"}
		case strings.Contains(err.Error(), "http: request body too large"):
			return &Err{Code: http.StatusRequestEntityTooLarge, message: "body too large"}
		case strings.Contains(err.Error(), "json: unknown field"):
			return &Err{Code: http.StatusBadRequest, message: fmt.Sprintf("contains unknown field %s", getUnknownField(err))}
		case errors.As(err, &syntaxError):
			return &Err{Code: http.StatusUnprocessableEntity, message: fmt.Sprintf("body contains invalid syntax at character %d", syntaxError.Offset)}
		case errors.As(err, &unmarshalError):
			return &Err{Code: http.StatusBadRequest, message: getTypeError(err)}
		default:
			return &Err{Code: http.StatusInternalServerError, message: err.Error()}
		}
	}

	err = decoder.Decode(&struct{}{})
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return &Err{Code: http.StatusBadRequest, message: "body contains multiple json objects"}
		}
	}

	return nil
}
