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

func TestWriteWithoutIndentation(t *testing.T) {
	j, w, _ := arrangeTest(t, "")
	header := make(http.Header)
	header.Set("Connection", "Keep-Alive")

	err := j.Write(w, http.StatusOK, Envelope{"name": "Jason"}, header)
	require.Nil(t, err)

	assert.Equal(t, ContentTypeJSON, w.Header().Get(ContentType))
}

func TestWriteWithIndentation(t *testing.T) {
	_, w, _ := arrangeTest(t, "")

	j := New(5, true, true)

	err := j.Write(w, http.StatusOK, Envelope{"name": "Jason"}, nil)
	require.Nil(t, err)

	assert.Equal(t, ContentTypeJSON, w.Header().Get(ContentType))
}

func TestWriteWithoutIndentationError(t *testing.T) {
	j, w, _ := arrangeTest(t, "")

	chanErr := make(chan int)

	err := j.Write(w, http.StatusOK, Envelope{"name": chanErr}, nil)
	require.NotNil(t, err)
}

func TestWriteWithIndentationError(t *testing.T) {
	_, w, _ := arrangeTest(t, "")

	j := New(5, true, true)

	chanErr := make(chan int)

	err := j.Write(w, http.StatusOK, Envelope{"name": chanErr}, nil)
	require.NotNil(t, err)
}
