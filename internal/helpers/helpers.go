package helpers

import (
	// "fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

// jsonResponse is the type used for generic JSON responses
type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func  ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

func  WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {

	var output []byte

		out, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return err
		}
		output=out



	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(output)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func  errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var customErr error

	switch {
	case strings.Contains(err.Error(), "SQLSTATE 23505"):
		customErr = errors.New("duplicate value violates unique constraint")
		statusCode = http.StatusForbidden
	case strings.Contains(err.Error(), "SQLSTATE 22001"):
		customErr = errors.New("the value you are trying to insert is too large")
		statusCode = http.StatusForbidden
	case strings.Contains(err.Error(), "SQLSTATE 23503"):
		customErr = errors.New("foreign key violation")
		statusCode = http.StatusForbidden
	default:
		customErr = err
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = customErr.Error()

	WriteJSON(w, statusCode, payload)
}

