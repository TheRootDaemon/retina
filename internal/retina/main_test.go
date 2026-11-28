package retina

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

//go:embed scripts/retina.js
var retinaScrapeScript string

//go:embed fixtures/retina_search_results.html
var retinaSearchResults string

var retinaBaseURL string

type testSearchEngine struct{}

func (*testSearchEngine) Name() string { return "retina" }

func (*testSearchEngine) BaseURL() string { return retinaBaseURL }

func (*testSearchEngine) SearchURL(query string) string {
	params := url.Values{
		"q": {query},
	}

	return retinaBaseURL + "?" + params.Encode()
}

func (*testSearchEngine) ScrapeScript() string {
	return evaluateModule(retinaScrapeScript)
}

var testServer *httptest.Server

func TestMain(m *testing.M) {
	testServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(retinaSearchResults))
		},
	))

	retinaBaseURL = testServer.URL

	code := m.Run()
	testServer.Close()
	os.Exit(code)
}

func setupRetina(t *testing.T) *Retina {
	t.Helper()
	r := New()
	r.SetupBrowser(context.Background())
	t.Cleanup(r.Close)
	return r
}
