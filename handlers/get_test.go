package handlers_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Mihai22125/URLShortenerAPI/data"
	"github.com/Mihai22125/URLShortenerAPI/handlers"
	"github.com/gorilla/mux"
)

// db act as a dummy test package level database.
var db data.Urls

// init initialize a dummy db with some data
func Init() {
	list := []*data.URL{
		{ID: 0, OriginalURL: "https://www.example1.com", ShortURL: "aaaaaaaa"},
		{ID: 1, OriginalURL: "https://www.example2.org", ShortURL: "abcdabcd"},
	}
	db.Init(list)
}

func TestGetURL(t *testing.T) {

	tt := []struct {
		name         string
		value        string
		err          error
		expectedCode int
	}{
		{name: "ShouldPass", value: "aaaaaaaa", err: nil, expectedCode: http.StatusMovedPermanently},
		{name: "ShouldFail", value: "bbbbbbbb", err: handlers.ErrShortenedNotExist, expectedCode: http.StatusNotFound},
	}

	// initialise dummy data slice
	Init()
	l := log.New(os.Stdout, "url-api-test", log.LstdFlags)
	dummyHandler := handlers.NewUrls(l, db)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/"+tc.value, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			vars := map[string]string{
				"shortURL": tc.value,
			}

			req = mux.SetURLVars(req, vars)

			rec := httptest.NewRecorder()
			dummyHandler.GetURL(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.err != nil {
				// do something
				if res.StatusCode != http.StatusNotFound {
					t.Errorf("expected status Not Found; got %v", res.StatusCode)
				}
				if msg := string(bytes.TrimSpace(b)); msg != tc.err.Error() {
					t.Errorf("expected message %q; got %q", tc.err, msg)
				}
				return
			}
			fmt.Println(string(b))

			if res.StatusCode != http.StatusMovedPermanently {
				t.Errorf("expected status Status Moved Permanently; got %v", res.Status)
			}
		})
	}
}
