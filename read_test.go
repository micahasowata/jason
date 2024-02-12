package jason

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMalformedRequest(t *testing.T) {
	// Arrange
	m := malformedRequest{
		status: http.StatusBadRequest,
		msg:    http.StatusText(http.StatusBadRequest),
	}

	// Act
	msg := m.Error()

	// Assert
	assert.Equal(t, msg, m.msg)
}
