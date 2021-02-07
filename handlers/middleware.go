package handlers

import (
	"context"
	"net/http"

	"github.com/Mihai22125/URLShortenerAPI/data"
)

// MiddlewareValidateURL validates the url in the request and calls next if ok
func (u *Urls) MiddlewareValidateURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		newURL := data.URL{}

		err := newURL.FromJSON(r.Body)

		if err != nil {
			u.l.Println("[MiddlewareValidateURL] error deserializing request body: ", err)
			http.Error(
				rw,
				ErrDeserializingReq.Error(),
				http.StatusUnprocessableEntity,
			)
			return
		}

		err = newURL.Validate()
		if err != nil {
			u.l.Println("[MiddlewareValidateURL] error validating request body: ", err)
			http.Error(
				rw,
				ErrInvalidReqBody.Error(),
				http.StatusBadRequest,
			)
			return
		}

		if data.ValidateURL(newURL.OriginalURL) == false {
			u.l.Println("[MiddlewareValidateURL] error invalid URL")
			http.Error(
				rw,
				ErrInvalidURL.Error(),
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
