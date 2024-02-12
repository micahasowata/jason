package jason

import jsoniter "github.com/json-iterator/go"

// Jason holds configuration settings for parsing and sending objects
type Jason struct {
	MaxBodySize           int64
	IndentResponse        bool
	DisallowUnknownFields bool
	parser                jsoniter.API
}

// New constructs an instance of Jason with the provided settings
func New(maxBodySize int64, indentResponse bool, disallowUnknownFields bool) *Jason {
	return &Jason{
		MaxBodySize:           maxBodySize,
		IndentResponse:        indentResponse,
		DisallowUnknownFields: disallowUnknownFields,
		parser:                jsoniter.ConfigCompatibleWithStandardLibrary,
	}
}
