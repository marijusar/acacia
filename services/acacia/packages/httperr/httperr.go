package httperr

import (
	"encoding/json"
	"errors"
	"net/http"
)

type HttpErrorResponse struct {
	Message string `json:"message"`
}

type InternalStatusError struct {
	error
	status int
}

func (e InternalStatusError) unwrap() error   { return e.error }
func (e InternalStatusError) httpStatus() int { return e.status }

func MarshalError(e error) []byte {
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
	return InternalStatusError{
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
			w.Write(MarshalError(internalError.unwrap()))
		}
	}
}

// WithMiddlewareErrorHandler wraps a middleware function that returns an error
// The middleware function receives the ResponseWriter, next handler, and request, and returns an error
// This is useful for middleware that needs to perform validation or authorization
func WithMiddlewareErrorHandler(middleware func(http.ResponseWriter, *http.Request, http.Handler) error) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := middleware(w, r, next)

			// If no error, the middleware already called next.ServeHTTP
			if err == nil {
				return
			}

			// Handle error
			w.Header().Set("Content-Type", "application/json")

			var internalError interface {
				unwrap() error
				httpStatus() int
			}

			if !errors.As(err, &internalError) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(MarshalError(errors.New("Internal server error")))
				return
			}

			if internalError.unwrap() != nil {
				w.WriteHeader(internalError.httpStatus())
				w.Write(MarshalError(internalError.unwrap()))
			}
		})
	}
}
