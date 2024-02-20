// Copyright (c) 2024 Micah Asowata

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package jason

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Err represents an error with an associated code and message
type Err struct {
	// Code is the recommended error code for error handling. It represents an HTTP status code,
	// providing guidance on how to handle the error, such as returning it to the client or logging it.
	Code int `json:"code"`

	// Msg contains the error message, providing detailed information about the encountered error.
	Msg string `json:"msg"`
}

// Error implements the error interface for Err, returning its message.
func (e *Err) Error() string {
	return e.Msg
}

// Read reads JSON data from the request body and decodes it into the provided destination interface.
// It checks the content type of the request and ensures that the destination is a pointer.
// Additionally, it sets limits on the request body size and handles various decoding errors.
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
		case j.disallowUnknownFields && isUnknownField(err):
			return &Err{Code: http.StatusBadRequest, Msg: fmt.Sprintf("request body contains unknown field %s", getFieldName(getStmt(err.Error())))}
		case containsImproperAssignment(err):
			return &Err{Code: http.StatusBadRequest, Msg: fmt.Sprintf("request body contains an invalid value for %s (at character %d)", getFieldName(getStmt(err.Error())), findErrorLocation(err.Error()))}
		case isBadlyFormedJSON(err):
			return &Err{Code: http.StatusBadRequest, Msg: fmt.Sprintf("request body contains badly formed JSON (at position %d)", findErrorLocation(err.Error()))}
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
