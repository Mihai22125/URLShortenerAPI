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

// MiddlewareValidateURL checks if given URL is valid
func (u *Urls) MiddlewareValidateURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		newURL := data.URL{}

		err := newURL.FromJSON(r.Body)

		if err != nil {
			u.l.Println("[ERROR] deserializing URL", err)
			http.Error(
				rw,
				"Error reading URL",
				http.StatusBadRequest,
			)
			return
		}

		err = newURL.Validate()
		if err != nil {
			u.l.Println("[ERROR] validating URL. URL is not valid", err)
			http.Error(
				rw,
				"URL is not valid",
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
