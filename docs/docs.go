// Package docs classification URL-shortener API
//
// Documentation of URL-shortener API
//
//     Schemes: http
//     BasePath: /api/v1/
//     Version: 1.0.0
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//     - string
//
//     Produces:
//     - string
//
//
// swagger:meta
package docs

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handlers

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body struct {
		Message string `json:"message"`
	}
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body struct {
		Messages []string `json:"messages"`
	}
}

// redirected succesfull
// swagger:response redirectResponse
type redirectResponseWrapper struct {
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// redirect error
// swagger:response redirectError
type redirectErrorWrapper struct {
	// shortened URL does not exist
	// in: body
	Body struct {
		// shortened URL doesn't exist
		// Required: true
		Message string
	}
}

// swagger:parameters redirectLink
type shortURLParam struct {
	// a string resulted from shortening an URL
	//
	// in: path
	// type: string
	// required: true
	// example: sdgh7e
	ShortURL string `json:"shortURL"`
}

// swagger:parameters addLink
type addLinkParam struct {
	// original link wanted to be shortened
	//
	// in: body
	// type: application/json
	// required: true
	// example: https://www.google.com
	Body struct {
		OriginalURL string `json:"url"`
	}
}

// Data structure representing a shortened link
// swagger:response addLinkSuccessResponse
type linkResponseWrapper struct {
	// Newly created shortened link
	// in: body
	// example: sdgh7e
	Body struct {
		ShortURL string `json:"shortURL"`
	}
}
