package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Mihai22125/URLShortenerAPI/data"
)

var db data.Urls

func TestAddURL(t *testing.T) {

	tt := []struct {
		name  string
		value string
	}{
		{name: "ShouldPass1", value: `{"url": "https://www.example.com"}`},
		{name: "ShouldPass2", value: `{"url": "https://www.example.com"}`},
	}

	// init initialize a dummy db with some data
	list := []*data.URL{
		{ID: 0, OriginalURL: "https://www.example1.com", ShortURL: "aaaaaaaa"},
		{ID: 1, OriginalURL: "https://www.example2.org", ShortURL: "abcdabcd"},
	}
	db.Init(list)

	// initialise dummy data slice
	l := log.New(os.Stdout, "url-api-test", log.LstdFlags)
	dummyHandler := NewUrls(l, db)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/v1/", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			reader := strings.NewReader(tc.value)
			newURL := data.URL{}
			err = newURL.FromJSON(reader)

			if err != nil {
				t.Fatalf("error deserializing request body: %v", err)
			}

			// add the product to the context
			ctx := context.WithValue(req.Context(), keyURL{}, newURL)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			dummyHandler.AddURL(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status Status OK; got %v", res.Status)
			}
			if msg := string(bytes.TrimSpace(b)); msg == "" {
				t.Errorf("got empty string as response message")
			}

		})
	}
}
