package retina

import (
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// baseOptions is the number of default allocator options.
var baseOptions = len(chromedp.DefaultExecAllocatorOptions) + 5

func TestRetina_allocateOptions(t *testing.T) {
	tests := []struct {
		name   string
		retina *Retina
		want   int
	}{
		{
			name:   "default flags",
			retina: New(),
			want:   baseOptions + len(defaultFlags),
		},
		{
			name:   "override flag",
			retina: New(WithFlags(map[string]any{"disable-dev-shm-usage": false})),
			want:   baseOptions + len(defaultFlags),
		},
		{
			name:   "custom flags",
			retina: New(WithFlags(map[string]any{"lang": "en-US", "mute-audio": false})),
			want:   baseOptions + len(defaultFlags) + 2,
		},
		{
			name:   "other options",
			retina: New(WithWindowSize(800, 600), WithUserAgent("retina/1.0"), WithTimeout(time.Second)),
			want:   baseOptions + len(defaultFlags),
		},
		{
			name:   "empty flags",
			retina: New(WithFlags(nil), WithFlags(map[string]any{})),
			want:   baseOptions + len(defaultFlags),
		},
		{
			name:   "zero value",
			retina: &Retina{},
			want:   baseOptions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.retina.allocateOptions()

			require.Len(t, got, tt.want)
		})
	}
}

func TestRetina_allocateOptions_PreservesDefaults(t *testing.T) {
	want := len(chromedp.DefaultExecAllocatorOptions)

	New(WithFlags(map[string]any{"lang": "en-US"})).allocateOptions()
	assert.Len(t, chromedp.DefaultExecAllocatorOptions, want)
}
