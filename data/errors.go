package data

import "errors"

// ErrURLNotFound Url not found in urlList
var ErrURLNotFound error = errors.New("URL not found")

// ErrFailedDecodeURL Failed to decode base62 short URL
var ErrFailedDecodeURL error = errors.New("Failed to decode base62 short URL")
