package main

import (
	"io/fs"
	"path/filepath"
	"text/template"
)

func newTemplateCache(uiFS fs.FS) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(uiFS, "pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		// Create a slice containing the filepaths for the base template and partials and page
		files := []string{
			"base.tmpl",
			"parts/nav.tmpl",
			"pages/" + name,
		}

		// Parse the files into template set
		ts, err := template.ParseFS(uiFS, files...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache map
		cache[name] = ts
	}

	return cache, nil
}
