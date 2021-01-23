package data

import (
	"encoding/json"
	"io"
	"net/url"

	"github.com/lytics/base62"
)

//URL struct for storing url details
type URL struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"url"`
	ShortURL    string `json:"shortened_url"`
}

// FromJSON extracts URL struct from JSON
func (u *URL) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
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

// EncodeURL shortens a long url using base62 encoding
func EncodeURL(longURL string) string {
	shortened := base62.StdEncoding.EncodeToString([]byte(longURL))
	return shortened
}

// DecodeURL returns original url from decoding base62 shortened url
func DecodeURL(shortURL string) (string, error) {
	decoded, err := base62.StdEncoding.DecodeString(shortURL)
	if err != nil {
		return "", ErrFailedDecodeURL
	}
	longURL := string(decoded)
	return longURL, nil
}

// ValidateURL checks if given URL is valid
func ValidateURL(givenURL string) error {
	_, err := url.ParseRequestURI(givenURL)
	return err
}
