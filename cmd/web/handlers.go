package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.data.Passphrase = " "
	app.data.Entropy = " "

	app.render(w, http.StatusOK, "home.tmpl", &app.data)
}

func (app *application) gen5Words(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Use the cached wordlist
	wordlist := app.wordlist

	// Generate a passphrase with 5 words (recommended for high security)
	wordCount := 5

	passphrase, err := generatePassphrase(wordlist, wordCount)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.data.Passphrase = passphrase
	app.data.Entropy = fmt.Sprintf("~%.1f Bits", float64(wordCount)*12.9)

	app.render(w, http.StatusOK, "home.tmpl", &app.data)

}
