package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type TestRequests []struct {
	name           string
	ignoreResponse bool
	response       string
	requestBody    string
	statusCode     int
}

func TestCreateEntry(t *testing.T) {
	tt := TestRequests{
		{
			name:           "body is nil",
			response:       "Not a valid URL",
			statusCode:     http.StatusBadRequest,
			ignoreResponse: true,
		},
		{
			name:        "short URL generation",
			requestBody: "https://www.google.de/",
			statusCode:  http.StatusOK,
		},
		{
			name:           "no http in url",
			requestBody:    "google.de",
			statusCode:     http.StatusOK,
			ignoreResponse: true,
		},
		{
			name:           "no valid URL",
			requestBody:    "this is really not a URL",
			statusCode:     http.StatusBadRequest,
			response:       "Not a valid URL",
			ignoreResponse: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			respBody := createRequestWithBody(t, tc.requestBody, tc.statusCode)
			if len(tc.response) > 0 {
				if string(respBody) != tc.response {
					t.Fatalf("expected body: %s; got: %s", tc.response, respBody)
				}
			}
			if tc.ignoreResponse {
				return
			}
			t.Run("test if shorted URL is correct", func(t *testing.T) {
				testRedirect(t, string(respBody), tc.requestBody)
			})
		})
	}
}

func TestSameUrl(t *testing.T) {
	t.Run("same url", func(t *testing.T) {

		respBody1 := createRequestWithBody(t, "yandex.ru", http.StatusOK)
		respBody2 := createRequestWithBody(t, "yandex.ru", http.StatusOK)
		if string(respBody1) != string(respBody2) {
			t.Fatalf("expected same body: %s; got: %s", respBody1, respBody2)
		}

		t.Run("test if shorted URL is correct", func(t *testing.T) {
			testRedirect(t, string(respBody1), "http://yandex.ru")
		})

	})
}

func createRequestWithBody(t *testing.T, requestBody string, statusCode int) []byte {
	resp, err := http.PostForm("http://0.0.0.0:8000/", url.Values{"url": {requestBody}})

	if err != nil {
		t.Fatalf("could not do request: %v", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("could not read body: %v", err)
	}
	if resp.StatusCode != statusCode {
		t.Errorf("expected status %d; got %d", statusCode, resp.StatusCode)
	}
	return bytes.TrimSpace(respBody)
}

func testRedirect(t *testing.T, shortURL, longURL string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}, // don't follow redirects
	}
	u, err := url.Parse(shortURL)
	if err != nil {
		t.Fatalf("could not parse shorted URL: %v", err)
	}
	resp, err := client.Do(&http.Request{
		URL: u,
	})
	if err != nil {
		t.Fatalf("could not do http request to shorted URL: %v", err)
	}
	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected status code: %d; got: %d", http.StatusFound, resp.StatusCode)
	}
	if resp.Header.Get("Location") != longURL {
		t.Fatalf("redirect URL is not correct. Expected: %s, got: %s", longURL, resp.Header.Get("Location"))
	}
}
