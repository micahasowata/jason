package v2

// Err is the underlying type of the error return from calling any Parser methods
type Err struct {
	Code    int
	message string
}

func (e *Err) Error() string {
	return e.message
}
