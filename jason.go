package jason

// Jason holds configuration settings for parsing and sending objects
type Jason struct {
	MaxBodySize           int64
	IndentResponse        bool
	DisallowUnknownFields bool
}

// New constructs an instance of Jason with the provided settings
func New(maxBodySize int64, indentResponse bool, disallowUnknownFields bool) *Jason {
	return &Jason{
		MaxBodySize:           maxBodySize,
		IndentResponse:        indentResponse,
		DisallowUnknownFields: disallowUnknownFields,
	}
}
