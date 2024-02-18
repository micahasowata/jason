package jason

import (
	"net/http"
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

func TestFindErrorLocation(t *testing.T) {
	errs := []string{
		"malformed body",
		"# body",
	}

	for _, err := range errs {
		got := findErrorLocation(err)

		assert.Equal(t, 0, got)
	}
}

func TestGetFieldName(t *testing.T) {
	input := "name=Jason"

	got := getFieldName(input)

	assert.Equal(t, "", got)
}