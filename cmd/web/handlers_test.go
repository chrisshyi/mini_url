package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGetURL(t *testing.T) {
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name            string
		urlPath         string
		wantStatusCode  int
		wantRedirectURL string
	}{
		{"Get existing short URL", "/b", 303, "http://mock.com"},
		{"Get non-existing short URL", "/a231", 404, ""},
		{"Get invalid short URL", "/你好ie1", 400, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, header, _ := ts.get(t, test.urlPath)

			if code != test.wantStatusCode {
				t.Errorf("want %d; got %d", test.wantStatusCode, code)
			}

			if !strings.Contains(header.Get("location"), test.wantRedirectURL) {
				t.Errorf("want body to contain %s", test.wantRedirectURL)
			}
		})
	}
}

func TestPostURL(t *testing.T) {
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name           string
		longURL        string
		wantStatusCode int
		wantShortURL   string
	}{
		{"Post existing URL", "http://mock.com", 200, "b"},
		{"Post non-existing URL", "http://www.google.com", 201, "c"},
		{"Post another non-existing URL", "http://www.gooogle.com", 201, "c"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			reqBody := NewURL{
				URL: test.longURL,
			}
			jsonBody, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatalf("Cannot marshal long URL %s", test.longURL)
			}
			code, _, respBodyJSON := ts.post(t, "/", jsonBody)

			var respBody ShortenedURL
			if code != test.wantStatusCode {
				t.Errorf("want %d; got %d", test.wantStatusCode, code)
			}
			// fmt.Printf("want short URL: %s, got %s", test.wantShortURL, string(respBodyJSON))

			err = json.Unmarshal(respBodyJSON, &respBody)
			if err != nil {
				t.Fatalf("Cannot unmarshal resp body for long URL %s", test.longURL)
			}
			if respBody.URL != test.wantShortURL {
				t.Errorf("want body to contain %s", test.wantShortURL)
			}
		})
	}

}
