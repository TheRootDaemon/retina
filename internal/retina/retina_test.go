package retina

import (
	"context"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// baseOptions is the number of default allocator options.
var baseOptions = len(chromedp.DefaultExecAllocatorOptions) + 5

func TestRetinaSetupBrowser(t *testing.T) {
	const timeout = 30 * time.Second

	r := New(
		WithUserAgent("MockRetina/1.0"),
		WithWindowSize(800, 600),
		WithTimeout(30*time.Second),
	)

	r.SetupBrowser(context.Background())
	t.Cleanup(func() { assert.NotPanics(t, r.Close) })

	require.NoError(t, chromedp.Do(r.ctx, chromedp.Navigate("about:blank")))

	t.Run("context not nil", func(t *testing.T) {
		require.NotNil(t, r.ctx)
		require.NotNil(t, r.cancel)
	})

	t.Run("context deadline", func(t *testing.T) {
		deadline, ok := r.ctx.Deadline()
		require.True(t, ok)
		assert.WithinDuration(
			t,
			time.Now().Add(30*time.Second),
			deadline,
			5*time.Second,
		)
	})

	t.Run("browser window size", func(t *testing.T) {
		got, err := chromedp.Run(
			r.ctx,
			chromedp.Evaluate[struct {
				Width  int `json:"width"`
				Height int `json:"height"`
			}](`({
				width: window.innerWidth,
				height: window.innerHeight,
			})`),
		)

		require.NoError(t, err)
		assert.Equal(t, 800, got.Width)
		assert.Equal(t, 600, got.Height)
	})

	t.Run("brower user agent", func(t *testing.T) {
		got, err := chromedp.Run(
			r.ctx,
			chromedp.Evaluate[struct {
				UserAgent string `json:"userAgent"`
			}](`({ userAgent: navigator.userAgent })`),
		)

		require.NoError(t, err)
		assert.Equal(t, "MockRetina/1.0", got.UserAgent)
	})
}

func TestRetinaClose(t *testing.T) {
	t.Run("close with setup", func(t *testing.T) {
		r := &Retina{
			width:     1920,
			height:    1080,
			timeout:   30 * time.Second,
			userAgent: "RetinaIntegrationTest/1.0",
		}

		r.SetupBrowser(context.Background())

		// context, cancellation should not be nil after initialization
		require.NotNil(t, r.ctx)
		require.NotNil(t, r.cancel)

		require.NotPanics(t, r.Close)

		// context, cancellation must be nil after cleaning up the browser
		require.Nil(t, r.ctx)
		require.Nil(t, r.cancel)
	})

	t.Run("close without setup", func(t *testing.T) {
		r := &Retina{}
		require.NotPanics(t, r.Close)
	})

	t.Run("close twice", func(t *testing.T) {
		r := &Retina{
			width:   1920,
			height:  1080,
			timeout: 30 * time.Second,
		}

		r.SetupBrowser(context.Background())

		require.NotPanics(t, r.Close)
		require.NotPanics(t, r.Close)
	})
}

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
