package handlers_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Mihai22125/URLShortenerAPI/data"
	"github.com/Mihai22125/URLShortenerAPI/handlers"
)

type middlewareHandler struct {
}

func (middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

func TestMiddlewareValidateURL(t *testing.T) {

	tt := []struct {
		name         string
		body         string
		err          error
		expectedCode int
	}{
		{name: "ShouldPass1", body: `{"url": "https://www.example.com"}`, err: nil},
		{name: "ShouldFail1", body: ``, err: handlers.ErrDeserializingReq, expectedCode: http.StatusUnprocessableEntity},
		{name: "ShouldFail2", body: `{"id": 1}`, err: handlers.ErrInvalidReqBody, expectedCode: http.StatusBadRequest},
		{name: "ShouldFail3", body: `{"url": " http://foo.com"}`, err: handlers.ErrInvalidURL, expectedCode: http.StatusBadRequest},
	}

	db := data.Urls{}
	l := log.New(os.Stdout, "url-api-test", log.LstdFlags)
	dummyHandler := handlers.NewUrls(l, db)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			testMiddle := middlewareHandler{}
			body := strings.NewReader(tc.body)
			req, err := http.NewRequest("POST", "/api/v1/", body)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			rec := httptest.NewRecorder()
			testMiddle.ServeHTTP(rec, req)

			mid := dummyHandler.MiddlewareValidateURL(testMiddle)
			mid.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if tc.err != nil {
				if res.StatusCode != tc.expectedCode {
					t.Errorf("expected status %v; got %v", tc.expectedCode, res.StatusCode)
				}
				if msg := string(bytes.TrimSpace(b)); msg != tc.err.Error() {
					t.Errorf("expected message %q; got %q", tc.err, msg)
				}
				return
			}

		})
	}
}
