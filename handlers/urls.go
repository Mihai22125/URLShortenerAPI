package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/Mihai22125/URLShortenerAPI/data"
)

// Urls is a http.Handler
type Urls struct {
	l *log.Logger
}

type keyURL struct{}

// NewUrls creates a Urls handler
func NewUrls(l *log.Logger) *Urls {
	return &Urls{l}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// MiddlewareValidateURL checks if given URL is valid
func (u *Urls) MiddlewareValidateURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		newURL := data.URL{}

		err := newURL.FromJSON(r.Body)

		if err != nil {
			u.l.Println("[MiddlewareValidateURL] error deserializing URL: ", err)
			http.Error(
				rw,
				"Error reading request body",
				http.StatusUnprocessableEntity,
			)
			return
		}

		err = newURL.Validate()
		if err != nil {
			u.l.Println("[MiddlewareValidateURL] error validating request body: ", err)
			http.Error(
				rw,
				"Error validating request body",
				http.StatusBadRequest,
			)
			return
		}

		if data.ValidateURL(newURL.OriginalURL) == false {
			u.l.Println("[MiddlewareValidateURL] error invalid URL")
			http.Error(
				rw,
				"Error invalid URL",
				http.StatusBadRequest,
			)
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), keyURL{}, newURL)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
