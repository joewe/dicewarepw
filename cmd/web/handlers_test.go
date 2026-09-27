package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"dicewarepw/ui"
)

func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache(ui.Files)
	if err != nil {
		t.Fatal(err)
	}

	wordlist, err := loadWordlist()
	if err != nil {
		t.Fatal(err)
	}

	return &application{
		infoLog:       log.New(io.Discard, "", 0),
		errorLog:      log.New(io.Discard, "", 0),
		templateCache: templateCache,
		wordlist:      wordlist,
		uiFS:          ui.Files,
	}
}

func TestHomeSecurityHeaders(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if got := resp.Header.Get("X-Frame-Options"); got != "DENY" {
		t.Errorf("expected X-Frame-Options DENY, got %q", got)
	}

	if got := resp.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("expected X-Content-Type-Options nosniff, got %q", got)
	}

	if got := resp.Header.Get("Referrer-Policy"); got != "strict-origin-when-cross-origin" {
		t.Errorf("expected Referrer-Policy strict-origin-when-cross-origin, got %q", got)
	}
}

func TestConcurrentGenerations(t *testing.T) {
	app := newTestApplication(t)
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	const concurrency = 30
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(count int) {
			defer wg.Done()

			words := "6"
			if count%2 == 0 {
				words = "5"
			}

			formData := url.Values{}
			formData.Set("words", words)

			resp, err := http.Post(ts.URL+"/generate", "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
			if err != nil {
				t.Errorf("request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("failed to read body: %v", err)
				return
			}

			if !strings.Contains(string(body), "Generierte Passphrase") {
				t.Errorf("expected generated passphrase in response")
			}
		}(i)
	}

	wg.Wait()
}
