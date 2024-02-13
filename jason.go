package jason

import jsoniter "github.com/json-iterator/go"

// Jason holds configuration settings for parsing and sending objects
type Jason struct {
	maxBodySize           int64
	indentResponse        bool
	disallowUnknownFields bool
	parser                jsoniter.API
}

// New constructs an instance of Jason with the provided settings
func New(maxBodySize int64, indentResponse bool, disallowUnknownFields bool) *Jason {
	return &Jason{
		maxBodySize:           maxBodySize,
		indentResponse:        indentResponse,
		disallowUnknownFields: disallowUnknownFields,
		parser:                jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}
