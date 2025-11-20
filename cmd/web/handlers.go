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
	app.data.Entropie = " "

	app.render(w, http.StatusOK, "home.tmpl", &app.data)
}

func (app *application) gen5Words(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// fmt.Println("Lade deutsche Diceware-Wortliste...")

	wordlist, err := loadWordlist()
	if err != nil {
		log.Fatalf("Fehler: %v", err)
	}

	// Generiere ein Passwort mit 5 Wörtern (empfohlen für hohe Sicherheit)
	wordCount := 5

	passphrase, err := generatePassphrase(wordlist, wordCount)
	if err != nil {
		log.Fatalf("Fehler beim Generieren der Passphrase: %v", err)
	}

	app.data.Passphrase = passphrase
	app.data.Entropie = fmt.Sprintf("~%.1f Bits", float64(wordCount)*12.9)

	app.render(w, http.StatusOK, "home.tmpl", &app.data)

}
