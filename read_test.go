package jason

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Jason"}`))
	require.Nil(t, err)
	r.Header.Set(ContentType, ContentTypeJSON)

	j := New(100, false, true)
	err = j.Read(w, r, &input)
	require.Nil(t, err)

	assert.Equal(t, "Jason", input.Name)
}
