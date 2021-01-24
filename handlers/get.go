package handlers

import (
	"net/http"

	"github.com/Mihai22125/URLShortenerAPI/data"
	"github.com/gorilla/mux"
)

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
	http.Redirect(rw, r, urlData.OriginalURL, http.StatusPermanentRedirect)
}
