package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"text/template"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	data          struct {
		Passphrase string
		Entropie   string
	}
}

const wordlistURL = "https://raw.githubusercontent.com/bjoernalbers/diceware-wordlist-german/refs/heads/main/wordlist-german-diceware.txt"

// loadWordlist lädt die Diceware-Wortliste von GitHub
func loadWordlist() (map[string]string, error) {
	resp, err := http.Get(wordlistURL)
	if err != nil {
		return nil, fmt.Errorf("fehler beim Laden der Wortliste: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP-Fehler: %s", resp.Status)
	}

	wordlist := make(map[string]string)
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 2 {
			code := parts[0]
			word := parts[1]
			wordlist[code] = word
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("fehler beim Parsen der Wortliste: %v", err)
	}

	return wordlist, nil
}

// rollDice simuliert einen Würfelwurf mit kryptographisch sicherem Zufall
func rollDice() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + 1, nil
}

// generateDicewareCode generiert einen 5-stelligen Code durch 5 Würfelwürfe
func generateDicewareCode() (string, error) {
	var code strings.Builder
	for i := 0; i < 5; i++ {
		roll, err := rollDice()
		if err != nil {
			return "", err
		}
		code.WriteString(fmt.Sprintf("%d", roll))
	}
	return code.String(), nil
}

// generatePassphrase generiert eine Diceware-Passphrase mit n Wörtern
func generatePassphrase(wordlist map[string]string, wordCount int) (string, error) {
	var words []string

	for i := 0; i < wordCount; i++ {
		code, err := generateDicewareCode()
		if err != nil {
			return "", fmt.Errorf("fehler beim Würfeln: %v", err)
		}

		word, exists := wordlist[code]
		if !exists {
			// Falls der Code nicht in der Liste ist, nochmal würfeln
			i--
			continue
		}

		words = append(words, word)
	}

	return strings.Join(words, "-"), nil
}

func main() {

	infoLog := log.New(log.Writer(), "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(log.Writer(), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     ":4000",
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starte Server auf %s", srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)

}
