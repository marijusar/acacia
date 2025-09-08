package httperr

import (
	"encoding/json"
	"errors"
	"net/http"
)

type HttpErrorResponse struct {
	Message string `json:"message"`
}

type internalStatusError struct {
	error
	status int
}

func (e internalStatusError) unwrap() error   { return e.error }
func (e internalStatusError) httpStatus() int { return e.status }

func marshalError(e error) []byte {
	s := HttpErrorResponse{
		Message: e.Error(),
	}

	b, err := json.Marshal(s)

	if err != nil {
		return nil
	}

	return b
}

func WithStatus(e error, status int) error {
	return internalStatusError{
		error:  e,
		status: status,
	}
}

func WithCustomErrorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		w.Header().Set("Content-Type", "application/json")

		var internalError interface {
			unwrap() error
			httpStatus() int
		}

		if !errors.As(err, &internalError) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if internalError.unwrap() != nil {
			w.WriteHeader(internalError.httpStatus())
			w.Write(marshalError(internalError.unwrap()))
		}
	}
}
