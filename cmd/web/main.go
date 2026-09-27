package main

import (
	"bufio"
	"crypto/rand"
	_ "embed"
	"fmt"
	"io/fs"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"dicewarepw/ui"
)

//go:embed wordlist-german-diceware.txt
var wordlistContent string

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	wordlist      map[string]string
	uiFS          fs.FS
}

// loadWordlist parses the embedded German Diceware wordlist
func loadWordlist() (map[string]string, error) {
	wordlist := make(map[string]string, 7776)
	scanner := bufio.NewScanner(strings.NewReader(wordlistContent))

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
		return nil, fmt.Errorf("error parsing wordlist: %v", err)
	}

	return wordlist, nil
}

// rollDice simulates a dice roll using cryptographically secure random number generation
func rollDice() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + 1, nil
}

// generateDicewareCode generates a 5-digit code by rolling 5 dice
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

// generatePassphrase generates a Diceware passphrase with n words
func generatePassphrase(wordlist map[string]string, wordCount int) (string, error) {
	var words []string

	for i := 0; i < wordCount; i++ {
		code, err := generateDicewareCode()
		if err != nil {
			return "", fmt.Errorf("error rolling dice: %v", err)
		}

		word, exists := wordlist[code]
		if !exists {
			// If the code is not in the list, roll again
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

	templateCache, err := newTemplateCache(ui.Files)
	if err != nil {
		errorLog.Fatal(err)
	}

	wordlist, err := loadWordlist()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: templateCache,
		wordlist:      wordlist,
		uiFS:          ui.Files,
	}

	port := "4000"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	addr := "0.0.0.0:" + port // uberspace needs 0.0.0.0

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
