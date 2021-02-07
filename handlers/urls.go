package handlers

import (
	"log"

	"github.com/Mihai22125/URLShortenerAPI/data"
)

// Urls is a http.Handler
type Urls struct {
	l       *log.Logger
	urlList data.Urls
}

type keyURL struct{}

// NewUrls creates a Urls handler
func NewUrls(l *log.Logger, urlList data.Urls) *Urls {
	return &Urls{l, urlList}
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}
