package retina

import (
	"github.com/chromedp/chromedp"
	"github.com/therootdaemon/retina/logger"
)

// SearchResult represents a single search result
// returned by the search engine.
type SearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

// Search searches for the given query
// and returns the results found by the browser.
//
// [Retina.SetupBrowser] must be called beforehand,
// or else it might use the zero values while searching
// leading to unwanted behaviour.
//
// An error is returned if the browser fails to navigate to the search URL
// or if the results cannot be extracted.
func (r *Retina) Search(
	engine SearchEngine,
	query string,
) ([]SearchResult, error) {
	url := engine.SearchURL(query)

	logger.Info(
		"retina: starting search with engine, query = %s, %s",
		engine.Name(),
		query,
	)

	if err := chromedp.Do(
		r.ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
	); err != nil {
		logger.Error("retina: failed to navigate search(%s): %v", url, err)
		return nil, err
	}

	results, err := chromedp.Run(
		r.ctx,
		chromedp.Evaluate[[]SearchResult](engine.ScrapeScript(), chromedp.EvalAwaitPromise),
	)
	if err != nil {
		logger.Error("retina: failed to extract results: %v", err)
		return nil, err
	}

	logger.Info(
		"retina: search completed: engine=%s, query=%s, len(results)=%d",
		engine.Name(),
		query,
		len(results),
	)

	return results, nil
}
