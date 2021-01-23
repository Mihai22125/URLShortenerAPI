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
	ID          int    `json:"id"`
	OriginalURL string `json:"url" validate:"required,url"`
	ShortURL    string `json:"shortened_url"`
}

// FromJSON extracts URL struct from JSON
func (u *URL) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// validate checks if given URL is valid
func (u *URL) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

// ValidateURL checks if given URL is valid
func validateURL(fl validator.FieldLevel) bool {
	u, err := url.Parse(fl.Field().String())
	if err != nil {
		return false
	}
	if u.Scheme == "" || u.Host == "" || u.Path == "" {
		return false
	}

	return true
}

// Urls is a collection of URL
type Urls []*URL

// urlList ia a list of Urls
var urlList = Urls{}

// AddURL adds a new URL to urlList
func AddURL(u *URL) {
	u.ID = getNextID()
	urlList = append(urlList, u)
}

// getNextID generates a new ID for a new entry in urlList
func getNextID() int {
	if len(urlList) == 0 {
		return 0
	}

	last := urlList[len(urlList)-1]
	return last.ID + 1
}

// GetURLByLong returns the URL struct with given original url from urlList
// returns ErrUrlNotFOund if url does not exists in urlList
func GetURLByLong(longURL string) (*URL, error) {
	for _, u := range urlList {
		if u.OriginalURL == longURL {
			return u, nil
		}
	}

	return nil, ErrURLNotFound
}

// GetURLByShort return the URL struct with given shortened URL from urlList
func GetURLByShort(shortURL string) (*URL, error) {
	for _, u := range urlList {
		if u.ShortURL == shortURL {
			return u, nil
		}
	}
	return nil, ErrURLNotFound
}

// ShortURL generates a random string with length 8 that is not present in database yet
func ShortURL(longURL string) string {
	shortened := ""
	for {
		shortened = randStringBytes(8)
		if _, err := GetURLByShort(shortened); err != nil {
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

/*
// DecodeURL returns original url from decoding base62 shortened url
func DecodeURL(shortURL string) (string, error) {
	decoded, err := base62.StdEncoding.DecodeString(shortURL)
	if err != nil {
		return "", ErrFailedDecodeURL
	}
	longURL := string(decoded)
	return longURL, nil
}
*/
