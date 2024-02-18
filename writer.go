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
	"net/http"
)

// Envelope is the jason equivalent of gin.H
type Envelope map[string]any

// Write formats the response and writes it to the client
func (j *Jason) Write(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	var msg []byte
	var err error

	switch j.indentResponse {
	case true:
		msg, err = j.parser.MarshalIndent(data, "", " ")
		if err != nil {
			return &Err{Code: http.StatusInternalServerError, Msg: err.Error()}
		}
	default:
		msg, err = j.parser.Marshal(data)
		if err != nil {
			return &Err{Code: http.StatusInternalServerError, Msg: err.Error()}
		}
	}

	msg = append(msg, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(msg)

	return nil
}
