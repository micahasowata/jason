package v2

const OneMB = 1048576

const ContentType = "application/json"

// Parser reads and writes any encoding/json compatible JSON values
type Parser struct {
	Indent bool
	Limit  int
}

// NewParser returns a new Parser that reads up to specified limit and optionally indent response
func NewParser(limit int, indent bool) *Parser {
	if limit < 1 {
		limit = 1
	}

	return &Parser{
		Limit:  limit,
		Indent: indent,
	}
}

// Default returns a Parser that can read up to 1 MB (1048576 bytes) of data and does not indent response
func Default() *Parser {
	return NewParser(OneMB, false)
}
