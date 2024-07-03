package v2

const (
	OneMB = 1048576
)

// Parser holds parsing configuration.
type Parser struct {
	Indent             bool
	AllowUnknownFields bool
	Limit              int
}

// NewParser creates a new Parser
func NewParser(limit int, indent bool, allowUnknownFields bool) *Parser {
	return &Parser{
		Limit:              limit,
		Indent:             indent,
		AllowUnknownFields: allowUnknownFields,
	}
}

// Default creates instantiates a Parser that can read up to 10MB (10485760 bytes) of data,
// doesn't indent response and disallows unknown fields
func Default() *Parser {
	return NewParser(OneMB*10, false, false)
}
