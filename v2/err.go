package v2

// Err implements the error interface
type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Err) Error() string {
	return e.Message
}
