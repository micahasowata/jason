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
