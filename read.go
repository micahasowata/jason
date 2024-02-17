package jason

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
)

var (
	ErrInvalidContentType = errors.New("jason: content type is not application/json")
)

// isBodyJSON ensures the media type of the request is set to application/json
func isBodyJSON(r *http.Request) bool {
	contentTypeHeader := r.Header.Get(ContentType)
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentTypeHeader, ";")[0]))
	return mediaType == ContentTypeJSON
}

// Read parses the provided request body, triages possible errors and formats their message
func (j *Jason) Read(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	ok := isBodyJSON(r)
	if !ok {
		return ErrInvalidContentType
	}

	r.Body = http.MaxBytesReader(w, r.Body, j.maxBodySize)

	dec := j.parser.NewDecoder(r.Body)

	if j.disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError
		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.Contains(err.Error(), "ReadObject: found"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown field %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
