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
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func arrangeTest(t *testing.T, body string) (*Jason, http.ResponseWriter, *http.Request) {
	t.Helper()
	j := New(100, false, true)

	w := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	require.Nil(t, err)

	req.Header.Set(ContentType, ContentTypeJSON)

	return j, w, req
}

func TestRead(t *testing.T) {
	var input struct {
		Name string `json:"name"`
	}

	j, w, r := arrangeTest(t, `{"name":"Jason"}`)

	err := j.Read(w, r, &input)
	require.Nil(t, err)

	assert.Equal(t, "Jason", input.Name)
}

func TestRequestBodyErrors(t *testing.T) {
	tests := []struct {
		Name string
		Body string
		Want string
	}{
		{
			Name: `syntax error`,
			Body: `<?xml version="1.0" encoding="UTF-8"?><note><to>Jason</to></note>`,
			Want: `request body contains badly formed JSON (at position 1)`,
		},
		{
			Name: `malformed body`,
			Body: `{"name": "Jason", }`,
			Want: `request body contains an invalid value for "name" (at character 10)`,
		},
		{
			Name: `invalid body`,
			Body: `["foo", "bar"]`,
			Want: `request body contains badly formed JSON (at position 1)`,
		},
		{
			Name: `invalid field value`,
			Body: `{"name": 123}'`,
			Want: `request body contains an invalid value for "name" (at character 10)`,
		},
		{
			Name: `empty body`,
			Body: ``,
			Want: `request body must not be empty`,
		},
		{
			Name: `double JSON objects`,
			Body: `{"name":"Jason"}{"name":"Xoxo"}`,
			Want: `request body must contain just one JSON object`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var input struct {
				Name string `json:"name"`
			}

			j, w, r := arrangeTest(t, tt.Body)
			err := j.Read(w, r, &input)
			require.NotNil(t, err)

			assert.Equal(t, tt.Want, err.Error())
		})
	}
}

func TestRequestBodyContentTypeChecker(t *testing.T) {
	var input struct {
		Name string `json:"name"`
	}

	j, w, r := arrangeTest(t, ``)
	r.Header.Set(ContentType, "text/html; charset=utf-8")

	err := j.Read(w, r, &input)
	require.NotNil(t, err)

	assert.Equal(t, "content type is not application/json", err.Error())
}

func TestRequestBodyDestinationChecker(t *testing.T) {
	var input struct {
		Name string `json:"name"`
	}

	j, w, r := arrangeTest(t, ``)

	err := j.Read(w, r, input)
	require.NotNil(t, err)

	assert.Equal(t, "destination is not a pointer", err.Error())
}

func TestRequestBodyUnknownFieldsChecker(t *testing.T) {
	var input struct {
		Name string `json:"name"`
	}

	_, w, r := arrangeTest(t, `{"another_name": "Xoxo"}`)
	j := New(100, false, true)

	err := j.Read(w, r, &input)
	require.NotNil(t, err)

	assert.Equal(t, `request body contains unknown field "another_name"`, err.Error())
}

func TestRequestBodySizeChecker(t *testing.T) {
	var input struct {
		Name string `json:"name"`
	}

	_, w, r := arrangeTest(t, `{"name": "XoxoJason"}`)
	j := New(5, false, true)

	err := j.Read(w, r, &input)
	require.NotNil(t, err)

	assert.Equal(t, "request body must not be larger than 5", err.Error())
}
