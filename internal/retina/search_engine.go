package retina

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
