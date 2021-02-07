package handlers

import "errors"

// ErrInvalidReqBody error validating request body
var ErrInvalidReqBody error = errors.New("Error validating request body")

// ErrInvalidURL URL given to be shortened is not valid
var ErrInvalidURL error = errors.New("invalid URL")

// ErrDeserializingReq Error deserializing request body
var ErrDeserializingReq = errors.New("Error reading request body")

// ErrShortenedNotExist shortened URL does not exist
var ErrShortenedNotExist = errors.New("shortened URL does not exist")
