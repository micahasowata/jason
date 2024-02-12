package jason_test

import (
	"testing"

	"github.com/micahasowata/jason"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	j := jason.Jason{
		MaxBodySize:           1,
		IndentResponse:        true,
		DisallowUnknownFields: true,
	}

	constructedJason := jason.New(1, true, true)

	require.EqualExportedValues(t, *constructedJason, j)
	assert.Equal(t, j.MaxBodySize, constructedJason.MaxBodySize)
	assert.Equal(t, j.IndentResponse, constructedJason.IndentResponse)
	assert.Equal(t, j.DisallowUnknownFields, constructedJason.DisallowUnknownFields)
}
