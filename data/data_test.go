package data

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {

	db := Urls{}

	list := []*URL{
		{ID: 0, OriginalURL: "https://www.example1.com", ShortURL: "aaaaaaaa"},
		{ID: 1, OriginalURL: "https://www.example2.org", ShortURL: "abcdabcd"},
	}
	db.Init(list)
}

func TestFromJSON(t *testing.T) {

	testStruct := &URL{}
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
		urlStruct URL
		expected  bool
	}{
		{
			"shouldPass1",
			URL{OriginalURL: "https://www.example.com"},
			true,
		},
		{
			"shouldPass2",
			URL{OriginalURL: "aaa"},
			true,
		},
		{
			"shouldFail1_empty_struct",
			URL{},
			false,
		},
		{
			"shouldFail2_no_url_field",
			URL{ID: 2},
			false,
		},
		{
			"shouldFail3_unexpected_field",
			URL{ID: 2, OriginalURL: "https://www.example.com"},
			false,
		},
		{
			"shouldFail4_unexpected_field",
			URL{ShortURL: "dsdf", OriginalURL: "https://www.example.com"},
			false,
		},
		{
			"shouldFail5_empty_url_string",
			URL{OriginalURL: ""},
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
			result := ValidateURL(tc.url)
			if result != tc.expected {
				t.Errorf("url: %v : expected %v; got %v", tc.url, tc.expected, result)
			}
		})
	}
}

func TestAddURL(t *testing.T) {

	testURLList := Urls{}
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
			testURLList.AddURL(&URL{OriginalURL: tc.longURL})
			urlFromSlice := testURLList.data[tc.expectedID]
			require.Equal(t, urlFromSlice.OriginalURL, tc.longURL)
			require.Equal(t, urlFromSlice.ID, tc.expectedID)
			require.NotNil(t, urlFromSlice.ShortURL)
		})
	}
}

func TestGetNextID(t *testing.T) {
	testURLList := Urls{}

	t.Run("empty database", func(t *testing.T) {
		id := testURLList.getNextID()
		require.Equal(t, id, 0)
	})
	testURLList.data = append(testURLList.data, &URL{ID: 0, OriginalURL: "http://www.example.com", ShortURL: "aaa"})

	t.Run("not empty database", func(t *testing.T) {
		id := testURLList.getNextID()
		require.Equal(t, id, 1)
	})

}

func TestGetURLByLong(t *testing.T) {
	testURLList := Urls{}
	tt := []struct {
		name          string
		longURL       string
		expectedError error
	}{
		{"shouldPass1", "https://www.example1.com", nil},
		{"shouldPass2", "https://www.example2.org", nil},
		{"shouldFail1", "http://www.sdfds.com", ErrURLNotFound},
	}

	// add some items in slice
	testURLList.data = append(testURLList.data, &URL{OriginalURL: tt[0].longURL})
	testURLList.data = append(testURLList.data, &URL{OriginalURL: tt[1].longURL})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			urlFromSlice, err := testURLList.GetURLByLong(tc.longURL)
			require.Equal(t, err, tc.expectedError)
			if err != nil {
				return
			}
			require.Equal(t, urlFromSlice.OriginalURL, tc.longURL)
		})
	}
}

func TestGetURLByShort(t *testing.T) {
	testURLList := Urls{}
	tt := []struct {
		name          string
		longURL       string
		shortURL      string
		expectedError error
	}{
		{"shouldPass1", "https://www.example1.com", "aaaaaa", nil},
		{"shouldPass2", "https://www.example2.org", "123456", nil},
		{"shouldFail1", "http://www.sdfds.com", "dvsdgf", ErrURLNotFound},
	}

	// add some items in slice
	testURLList.data = append(testURLList.data, &URL{OriginalURL: tt[0].longURL, ShortURL: tt[0].shortURL})
	testURLList.data = append(testURLList.data, &URL{OriginalURL: tt[1].longURL, ShortURL: tt[1].shortURL})

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			urlFromSlice, err := testURLList.GetURLByShort(tc.shortURL)
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
	testURLList := Urls{}
	tt := []struct {
		name          string
		longURL       string
		expectedError error
	}{
		{"shouldPass1", "https://www.example1.com", nil},
		{"shouldPass2", "https://www.example2.org", nil},
	}
	testURLList.data = append(testURLList.data, &URL{OriginalURL: tt[0].longURL})
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			shortURL := testURLList.ShortURL(tc.longURL)
			// require.Regexp(t, regexp.MustCompile("(?i)[a-z]{8}"), randStr)
			require.NotNil(t, shortURL)
		})
	}
}
