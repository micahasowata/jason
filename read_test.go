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
