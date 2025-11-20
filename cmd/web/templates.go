package main

import (
	"path/filepath"
	"text/template"
)

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// Create a slice containing the filepaths for the base template and partials and page
		files := []string{
			"./ui/base.tmpl",
			"./ui/parts/nav.tmpl",
			page,
		}

		// Parse the files into template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache map
		cache[name] = ts
	}

	return cache, nil
}
