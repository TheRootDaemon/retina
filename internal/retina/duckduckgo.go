package retina

import _ "embed"

import "net/url"

// duckDuckGoBaseURL is the base URL used for DuckDuckGo searches.
const duckDuckGoBaseURL = "https://html.duckduckgo.com/html/"

//go:embed duckduckgo.js
var duckDuckGoScrapeScript string

// DuckDuckGo represents the DuckDuckGo search engine,
// and implements [SearchEngine].
type DuckDuckGo struct{}

// Name returns the name of the search engine.
func (d *DuckDuckGo) Name() string { return "duckduckgo" }

// BaseURL returns the base URL of DuckDuckGo.
func (d *DuckDuckGo) BaseURL() string { return duckDuckGoBaseURL }

// ScrapeScript returns valid JavaScript code
// that extracts search results from DuckDuckGo's results page.
//
// The script evaluates to an array of [SearchResult] objects.
func (d *DuckDuckGo) ScrapeScript() string { return duckDuckGoScrapeScript }

// SearchURL returns a URL for searching the given query on DuckDuckGo.
func (d *DuckDuckGo) SearchURL(query string) string {
	params := url.Values{
		"q": {query},
	}

	return duckDuckGoBaseURL + "?" + params.Encode()
}
