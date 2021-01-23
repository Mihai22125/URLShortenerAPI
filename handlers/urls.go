package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/Mihai22125/URLShortenerAPI/data"
	"github.com/gorilla/mux"
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

// GetURL redirects shortened URL to original URL
func (u *Urls) GetURL(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET URL")

	vars := mux.Vars(r)
	shortURL := vars["shortURL"]

	// check if given shortened url exists in database
	urlData, err := data.GetURLByShort(shortURL)
	if err != nil {
		http.Error(rw, "shortened URL does not exist", http.StatusNotFound)
		return
	}
	// redirect short url to coresponding web page
	http.Redirect(rw, r, urlData.OriginalURL, http.StatusSeeOther)
}

// AddURL adds a new URL struct in urlList
func (u *Urls) AddURL(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST URl")

	newURL := r.Context().Value(keyURL{}).(data.URL)

	shortURL := data.EncodeURL(newURL.OriginalURL)
	newURL.ShortURL = shortURL

	data.AddURL(&newURL)

	// return shortened URL for given link
	rw.Write([]byte(newURL.ShortURL))
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
