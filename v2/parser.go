package v2

const OneMB = 1048576

const ContentType = "application/json"

// Parser holds parsing configuration.
type Parser struct {
	Indent bool
	Limit  int
}

// NewParser creates a new Parser
func NewParser(limit int, indent bool) *Parser {
	if limit < 1 {
		limit = 1
	}

	return &Parser{
		Limit:  limit,
		Indent: indent,
	}
}

// Default creates instantiates a Parser that can read up to 10MB (10485760 bytes) of data,
// doesn't indent response and disallows unknown fields
func Default() *Parser {
	return NewParser(OneMB, false)
}
