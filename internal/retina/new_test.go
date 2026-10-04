package retina

import (
	"maps"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		want    Retina
	}{
		{
			name: "applies defaults",
			want: Retina{
				height:    defaultHeight,
				width:     defaultWidth,
				timeout:   defaultTimeout,
				userAgent: defaultUserAgent,
				flags:     flagsWith(nil),
			},
		},
		{
			name: "applies every option",
			options: []Option{
				WithWindowSize(800, 600),
				WithTimeout(5 * time.Second),
				WithUserAgent("retina/1.0"),
				WithFlags(map[string]any{"lang": "en-US"}),
			},
			want: Retina{
				height:    600,
				width:     800,
				timeout:   5 * time.Second,
				userAgent: "retina/1.0",
				flags:     flagsWith(map[string]any{"lang": "en-US"}),
			},
		},
		{
			name: "applies the last option",
			options: []Option{
				WithWindowSize(100, 200),
				WithWindowSize(300, 400),
				WithTimeout(time.Second),
				WithTimeout(2 * time.Second),
				WithUserAgent("first"),
				WithUserAgent("second"),
				WithFlags(map[string]any{"lang": "en-US"}),
				WithFlags(map[string]any{"lang": "fr-FR"}),
			},
			want: Retina{
				height:    400,
				width:     300,
				timeout:   2 * time.Second,
				userAgent: "second",
				flags:     flagsWith(map[string]any{"lang": "fr-FR"}),
			},
		},
		{
			name: "applies zero values",
			options: []Option{
				WithWindowSize(0, 0),
				WithTimeout(0),
				WithUserAgent(""),
				WithFlags(nil),
				WithFlags(map[string]any{}),
			},
			want: Retina{flags: flagsWith(nil)},
		},
		{
			name: "keeps the default flags when overriding one",
			options: []Option{
				WithFlags(map[string]any{"disable-blink-features": "AutomationControlled2"}),
			},
			want: Retina{
				height:    defaultHeight,
				width:     defaultWidth,
				timeout:   defaultTimeout,
				userAgent: defaultUserAgent,
				flags:     flagsWith(map[string]any{"disable-blink-features": "AutomationControlled2"}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.options...)

			require.NotNil(t, got)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestNew_IsolatesFlags(t *testing.T) {
	defaults := New()

	flags := map[string]any{"lang": "en-US"}
	configured := New(WithFlags(flags))

	t.Run("does not retain the given map", func(t *testing.T) {
		flags["lang"] = "fr-FR"
		assert.Equal(t, "en-US", configured.flags["lang"])
	})

	t.Run("does not leak flags between instances", func(t *testing.T) {
		assert.Equal(t, "en-US", configured.flags["lang"])
		assert.NotContains(t, defaults.flags, "lang")
	})
}

// flagsWith returns [defaultFlags] merged with the given overrides.
func flagsWith(overrides map[string]any) map[string]any {
	flags := maps.Clone(defaultFlags)
	maps.Copy(flags, overrides)
	return flags
}
