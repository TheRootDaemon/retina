package retina

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// duckDuckGo is the [SearchEngine] that every test in this file exercises.
var duckDuckGo SearchEngine = &DuckDuckGo{}

func TestDuckDuckGo_Name(t *testing.T) {
	assert.Equal(t, "duckduckgo", duckDuckGo.Name())
}

func TestDuckDuckGo_BaseURL(t *testing.T) {
	assert.Equal(t, "https://html.duckduckgo.com/html/", duckDuckGo.BaseURL())
}

func TestDuckDuckGo_ScrapeScript(t *testing.T) {
	script := duckDuckGo.ScrapeScript()
	assert.NotEmpty(t, script)
}

func TestDuckDuckGo_SearchURL(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "single word",
			query: "golang",
			want:  "https://html.duckduckgo.com/html/?q=golang",
		},
		{
			name:  "words",
			query: "watchman hop retina",
			want:  "https://html.duckduckgo.com/html/?q=watchman+hop+retina",
		},
		{
			name:  "empty",
			query: "",
			want:  "https://html.duckduckgo.com/html/?q=",
		},
		{
			name:  "ampersand and equals",
			query: "a&b=c",
			want:  "https://html.duckduckgo.com/html/?q=a%26b%3Dc",
		},
		{
			name:  "plus",
			query: "c++",
			want:  "https://html.duckduckgo.com/html/?q=c%2B%2B",
		},
		{
			name:  "percent",
			query: "100%",
			want:  "https://html.duckduckgo.com/html/?q=100%25",
		},
		{
			name:  "slash",
			query: "/slash",
			want:  "https://html.duckduckgo.com/html/?q=%2Fslash",
		},
		{
			name:  "non ascii",
			query: "café",
			want:  "https://html.duckduckgo.com/html/?q=caf%C3%A9",
		},
		{
			name:  "tab",
			query: "a\tb",
			want:  "https://html.duckduckgo.com/html/?q=a%09b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := duckDuckGo.SearchURL(tt.query)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDuckDuckGo_SearchURL_RoundTrips(t *testing.T) {
	queries := []string{
		"retin",
		"watchman hop retina",
		"",
		"a&b=c",
		"c++",
		"100%",
		"/slash",
		"café",
		"a\tb",
		"?q=injected",
		"#fragment",
	}

	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			parsed, err := url.Parse(duckDuckGo.SearchURL(query))
			require.NoError(t, err)

			assert.True(t, strings.HasPrefix(parsed.String(), duckDuckGo.BaseURL()))
			assert.Equal(t, query, parsed.Query().Get("q"))
		})
	}
}

func TestDuckDuckGo_SearchURL_KeepsOneParameter(t *testing.T) {
	parsed, err := url.Parse(duckDuckGo.SearchURL("a&b=c"))
	require.NoError(t, err)

	assert.Len(t, parsed.Query(), 1)
}
