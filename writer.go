package jason

import (
	"net/http"
)

// Envelope is the jason equivalent of gin.H
type Envelope map[string]any

// Write formats the response and writes it to the client
func (j Jason) Write(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	var msg []byte
	var err error

	switch {
	case j.IndentResponse:
		msg, err = j.parser.MarshalIndent(data, "", "\t")
		if err != nil {
			return err
		}
	default:
		msg, err = j.parser.Marshal(data)
		if err != nil {
			return err
		}
	}

	msg = append(msg, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(msg)

	return nil
}
