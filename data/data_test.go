package data_test

import (
	"io"
	"strings"
	"testing"

	"github.com/Mihai22125/URLShortenerAPI/data"
	"github.com/stretchr/testify/require"
)

func TestFromJSON(t *testing.T) {

	testStruct := &data.URL{}
	tt := []struct {
		name     string
		body     io.Reader
		expected bool
	}{
		{
			"shouldPass1",
			strings.NewReader("{}"),
			true,
		},
		{
			"shouldPass2",
			strings.NewReader(`{"id":1, "url": "some", "shortened_url": "x"}`),
			true,
		},
		{
			"shouldFail1",
			strings.NewReader(("")),
			false,
		},
		{
			"shouldFail2",
			strings.NewReader("{"),
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := testStruct.FromJSON(tc.body)
			if (result == nil) != tc.expected {
				t.Errorf("expected %v; got %v", tc.expected, result)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	// tt - testing table
	tt := []struct {
		name      string
		urlStruct data.URL
		expected  bool
	}{
		{
			"shouldPass1",
			data.URL{OriginalURL: "https://www.example.com"},
			true,
		},
		{
			"shouldPass2",
			data.URL{OriginalURL: "aaa"},
			true,
		},
		{
			"shouldFail1_empty_struct",
			data.URL{},
			false,
		},
		{
			"shouldFail2_no_url_field",
			data.URL{ID: 2},
			false,
		},
		{
			"shouldFail3_unexpected_field",
			data.URL{ID: 2, OriginalURL: "https://www.example.com"},
			false,
		},
		{
			"shouldFail4_unexpected_field",
			data.URL{ShortURL: "dsdf", OriginalURL: "https://www.example.com"},
			false,
		},
		{
			"shouldFail5_empty_url_string",
			data.URL{OriginalURL: ""},
			false,
		},
	}

	// tc - test case
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.urlStruct.Validate()
			if (result == nil) != tc.expected {
				t.Errorf("expected %v; got %v", tc.expected, result)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tt := []struct {
		name     string
		url      string
		expected bool
	}{
		{"shouldPass1", "https://www.example.com", true},
		{"shouldPass2", "http://www.example.com", true},
		{"shouldPass3", "http://www.example.x/something", true},
		{"shouldPass4", "http://www.example", true},

		{"shouldFail1", "www.example.com", false},
		{"shouldFail2", "example.com", false},
		{"shouldFail3", "randomtext", false},
		{"shouldFail4", "www. example.com", false},
		{"shouldFail5", " http://foo.com", false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := data.ValidateURL(tc.url)
			if result != tc.expected {
				t.Errorf("url: %v : expected %v; got %v", tc.url, tc.expected, result)
			}
		})
	}
}

func TestAddURL(t *testing.T) {

	testUrlList := data.Urls{}
	tt := []struct {
		name       string
		longURL    string
		expectedID int
	}{
		{"adding item in empty slice", "https://www.example1.com", 0},
		{"adding item in not empty slice", "https://www.example2.org", 1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testUrlList.AddURL(&data.URL{OriginalURL: tc.longURL})
			urlFromSlice, err := testUrlList.GetURLByLong(tc.longURL)
			require.Equal(t, err, nil)
			require.Equal(t, urlFromSlice.OriginalURL, tc.longURL)
			require.Equal(t, urlFromSlice.ID, tc.expectedID)
			require.NotNil(t, urlFromSlice.ShortURL)
		})
	}
}

func TestGetNextID(t *testing.T) {
	testUrlList := data.Urls{}
	tt := []struct {
		name       string
		longURL    string
		expectedID int
	}{
		{"first item ID in slice", "https://www.example1.com", 0},
		{"second item ID slice", "https://www.example2.org", 1},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testUrlList.AddURL(&data.URL{OriginalURL: tc.longURL})
			urlFromSlice, err := testUrlList.GetURLByLong(tc.longURL)
			require.Equal(t, err, nil)
			require.Equal(t, urlFromSlice.ID, tc.expectedID)
		})
	}
}

func TestGetURLByLong(t *testing.T) {
	testUrlList := data.Urls{}
	tt := []struct {
		name          string
		longURL       string
		expectedError error
	}{
		{"shouldPass1", "https://www.example1.com", nil},
		{"shouldPass2", "https://www.example2.org", nil},
		{"shouldFail1", "http://www.sdfds.com", data.ErrURLNotFound},
	}

	// add some items in slice
	testUrlList.AddURL(&data.URL{OriginalURL: tt[0].longURL})
	testUrlList.AddURL(&data.URL{OriginalURL: tt[1].longURL})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			urlFromSlice, err := testUrlList.GetURLByLong(tc.longURL)
			require.Equal(t, err, tc.expectedError)
			if err != nil {
				return
			}
			require.Equal(t, urlFromSlice.OriginalURL, tc.longURL)
		})
	}
}

func TestGetURLByShort(t *testing.T) {
	testUrlList := data.Urls{}
	tt := []struct {
		name          string
		longURL       string
		shortURL      string
		expectedError error
	}{
		{"shouldPass1", "https://www.example1.com", "aaaaaa", nil},
		{"shouldPass2", "https://www.example2.org", "123456", nil},
		{"shouldFail1", "http://www.sdfds.com", "dvsdgf", data.ErrURLNotFound},
	}

	// add some items in slice
	testUrlList.AddURL(&data.URL{OriginalURL: tt[0].longURL, ShortURL: tt[0].shortURL})
	testUrlList.AddURL(&data.URL{OriginalURL: tt[1].longURL, ShortURL: tt[1].shortURL})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			urlFromSlice, err := testUrlList.GetURLByShort(tc.shortURL)
			require.Equal(t, err, tc.expectedError)
			if err != nil {
				return
			}
			require.Equal(t, urlFromSlice.OriginalURL, tc.longURL)
			require.Equal(t, urlFromSlice.ShortURL, tc.shortURL)
		})
	}
}

func TestShortURL(t *testing.T) {
	testUrlList := data.Urls{}
	tt := []struct {
		name          string
		longURL       string
		expectedError error
	}{
		{"shouldPass1", "https://www.example1.com", nil},
		{"shouldPass2", "https://www.example2.org", nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			shortURL := testUrlList.ShortURL(tc.longURL)
			// require.Regexp(t, regexp.MustCompile("(?i)[a-z]{8}"), randStr)
			require.NotNil(t, shortURL)
			testUrlList.AddURL(&data.URL{OriginalURL: tt[0].longURL, ShortURL: shortURL})
		})
	}
}
