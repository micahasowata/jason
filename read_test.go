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

func TestIsBodyJSON(t *testing.T) {
	tests := []struct {
		Name          string
		AddJSONHeader bool
		Want          bool
	}{
		{
			Name:          "valid content type header",
			AddJSONHeader: true,
			Want:          true,
		},
		{
			Name:          "invalid content type header",
			AddJSONHeader: false,
			Want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			r, err := http.NewRequest(http.MethodPost, "/", nil)
			require.Nil(t, err)

			if tt.AddJSONHeader {
				r.Header.Set(ContentType, ContentTypeJSON)
			}

			got := isBodyJSON(r)

			assert.Equal(t, tt.Want, got)
		})
	}
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

func TestReadErrors(t *testing.T) {
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

func TestContentType(t *testing.T) {
	var input struct {
		Name string `json:"name"`
	}

	j, w, r := arrangeTest(t, ``)
	r.Header.Set(ContentType, "text/html; charset=utf-8")

	err := j.Read(w, r, &input)
	require.NotNil(t, err)

	assert.Equal(t, "content type is not application/json", err.Error())
}

func TestDestinationCheck(t *testing.T) {
	var input struct {
		Name string `json:"name"`
	}

	j, w, r := arrangeTest(t, ``)

	err := j.Read(w, r, input)
	require.NotNil(t, err)

	assert.Equal(t, "destination is not a pointer", err.Error())
}
