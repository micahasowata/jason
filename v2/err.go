package v2

// Err implements the error interface
type Err struct {
	Code    int
	message string
}

func (e *Err) Error() string {
	return e.message
}
