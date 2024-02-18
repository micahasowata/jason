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
