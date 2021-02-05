package data

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/url"
	"time"

	"github.com/go-playground/validator"
)

//URL struct for storing url details
type URL struct {
	ID          int    `json:"id" validate:"isdefault"`
	OriginalURL string `json:"url" validate:"required"`
	ShortURL    string `json:"shortened_url" validate:"isdefault"`
}

// FromJSON extracts URL struct from JSON
func (u *URL) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// Validate checks if URL struct valid
func (u *URL) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// ValidateURL checks if given URL is valid
func ValidateURL(longURL string) bool {
	u, err := url.Parse(longURL)
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

// Urls is a collection of URL
type Urls struct {
	data []*URL
}

// urlList ia a list of Urls
var urlList = Urls{}

// AddURL adds a new URL to urlList
func (urlList *Urls) AddURL(u *URL) {
	u.ID = urlList.getNextID()
	urlList.data = append(urlList.data, u)
}

// getNextID generates a new ID for a new entry in urlList
func (urlList *Urls) getNextID() int {
	if len(urlList.data) == 0 {
		return 0
	}

	last := urlList.data[len(urlList.data)-1]
	return last.ID + 1
}

// GetURLByLong returns the URL struct with given original url from urlList
// returns ErrUrlNotFOund if url does not exists in urlList
func (urlList *Urls) GetURLByLong(longURL string) (*URL, error) {
	for _, u := range urlList.data {
		if u.OriginalURL == longURL {
			return u, nil
		}
	}

	return nil, ErrURLNotFound
}

// GetURLByShort return the URL struct with given shortened URL from urlList
func (urlList *Urls) GetURLByShort(shortURL string) (*URL, error) {
	for _, u := range urlList.data {
		if u.ShortURL == shortURL {
			return u, nil
		}
	}
	return nil, ErrURLNotFound
}

// ShortURL generates a random string with length 8 that is not present in database yet
func (urlList *Urls) ShortURL(longURL string) string {
	shortened := ""
	for {
		shortened = randStringBytes(8)
		if _, err := urlList.GetURLByShort(shortened); err != nil {
			break
		}
	}

	return shortened
}

func randStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
