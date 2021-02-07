package handlers

import (
	"net/http"

	"github.com/Mihai22125/URLShortenerAPI/data"
)

// swagger:route POST /  postLink addLink
// adds new link to datebase and returns respective shortened string
// consumes:
//     -application/json
//
// produces:
//	   -string
// responses:
//	200: addLinkSuccessResponse
//  422: errorValidation
//  400: errorResponse

// AddURL adds a new URL struct in urlList
func (u *Urls) AddURL(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle POST URl")

	newURL := r.Context().Value(keyURL{}).(data.URL)

	shortURL := u.urlList.ShortURL(newURL.OriginalURL)
	newURL.ShortURL = shortURL

	u.urlList.AddURL(&newURL)

	// return shortened URL for given link
	rw.Write([]byte(newURL.ShortURL))

}
