package jason_test

import (
	"testing"

	"github.com/micahasowata/jason"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	bodySize := 1
	indentResponse := true
	disallowUnknownFields := true

	constructedJason := jason.New(1, true, true)

	assert.Equal(t, int64(bodySize), constructedJason.MaxBodySize)
	assert.Equal(t, indentResponse, constructedJason.IndentResponse)
	assert.Equal(t, disallowUnknownFields, constructedJason.DisallowUnknownFields)
}
