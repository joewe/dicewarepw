package main

import (
	"io/fs"
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/generate", app.generate)

	staticFS, err := fs.Sub(app.uiFS, "static")
	if err != nil {
		app.errorLog.Println(err)
	}

	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
