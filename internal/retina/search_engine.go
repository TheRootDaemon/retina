package retina

import (
	"encoding/base64"
	"fmt"
)

// SearchEngine represents a search engine supported by [Retina].
type SearchEngine interface {
	// Name returns the name of the search engine.
	Name() string

	// BaseURL returns the base URL of the search engine.
	BaseURL() string

	// ScrapeScript returns the JavaScript code,
	// used to extract the search results
	// from the search engine's page.
	ScrapeScript() string

	// SearchURL returns the query encoded URL,
	// according to the search engine's URL requirements.
	SearchURL(query string) string
}

// evaluateModule returns JavaScript
// that dynamically imports the given JavaScript module
// and invokes its default export.
//
// The module is encoded as a data URL
// so it can be evaluated directly
// in the browser context.
func evaluateModule(module string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(module))

	return fmt.Sprintf(`
		import("data:text/javascript;base64,%s")
			.then(module => module.default())
	`, encoded)
}
