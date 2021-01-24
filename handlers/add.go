package handlers

import (
	"net/http"

	"github.com/Mihai22125/URLShortenerAPI/data"
)

// AddURL adds a new URL struct in urlList
func (u *Urls) AddURL(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST URl")

	newURL := r.Context().Value(keyURL{}).(data.URL)

	shortURL := data.ShortURL(newURL.OriginalURL)
	newURL.ShortURL = shortURL

	data.AddURL(&newURL)

	// return shortened URL for given link
	rw.Write([]byte(newURL.ShortURL))
}