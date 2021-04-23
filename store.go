package embedqueries

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

//go:embed queries/*
var embeddedFS embed.FS

// newStore loads and returns all available queries
// Queries are loaded from the queries sub dir. See above, where it is declared to embed this dir
// returned map:
// 	key is a query identifier (the file name minus extension)
// 	value is the raw query (file contents), stored as a parsed Go Template to support easy variable substitution
func newStore() (map[string]*template.Template, error) {

	store := make(map[string]*template.Template)

	// load all embedded files into map
	entries, err := embeddedFS.ReadDir("queries")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read directory")
	}

	for _, entry := range entries {

		fileNameWithoutExtension := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

		b, err := embeddedFS.ReadFile(fmt.Sprintf("queries/%s", entry.Name()))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Failed to read file %s", entry.Name()))
		}

		tmpl := template.New(fileNameWithoutExtension)
		tmpl, err = tmpl.Parse(string(b))
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Failed to parse file %s", entry.Name()))
		}

		store[fileNameWithoutExtension] = tmpl
	}

	return store, nil
}
