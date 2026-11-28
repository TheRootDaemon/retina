package retina

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubSearchEngine struct {
	url    string
	script string
}

func (*stubSearchEngine) Name() string {
	return "stub"
}

func (e *stubSearchEngine) BaseURL() string {
	return e.url
}

func (e *stubSearchEngine) SearchURL(string) string {
	return e.url
}

func (e *stubSearchEngine) ScrapeScript() string {
	return e.script
}

func TestRetinaSearch_Success(t *testing.T) {
	engine := &testSearchEngine{}

	r := setupRetina(t)
	results, err := r.Search(engine, "retina")

	assert.Nil(t, err)
	assert.Equal(t, 3, len(results))
}

func TestRetinaSearch_InvalidResultType(t *testing.T) {
	r := setupRetina(t)

	engine := &stubSearchEngine{
		url: retinaBaseURL,
		script: `export default function scrape() {
			({ message: "not an array" })
		}`,
	}

	results, err := r.Search(engine, "retina")
	assert.Error(t, err)
	assert.Nil(t, results)
}

func TestRetinaSearch_NavigationError(t *testing.T) {
	r := setupRetina(t)

	engine := &stubSearchEngine{
		url:    "http://127.0.0.1:1/",
		script: "[]",
	}

	results, err := r.Search(engine, "retina")
	assert.Error(t, err)
	assert.Nil(t, results)
}

func TestRetinaSearch_ExtractionError(t *testing.T) {
	r := setupRetina(t)

	engine := &stubSearchEngine{
		url:    testServer.URL,
		script: `(() => { throw new Error("scraper failed"); })()`,
	}

	results, err := r.Search(engine, "retina")
	assert.Error(t, err)
	assert.Nil(t, results)
}
