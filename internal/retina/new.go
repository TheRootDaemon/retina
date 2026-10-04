package retina

import (
	"maps"
	"time"
)

const (
	// defaultHeight is the default browser viewport height used by [Retina].
	defaultHeight = 1080

	// defaultWidth is the default browser viewport width used by [Retina].
	defaultWidth = 1920

	// defaultTimeout is the default maximum duration for a [Retina] browser context.
	defaultTimeout = 120 * time.Second

	// defaultUserAgent is the default user agent string used by [Retina].
	defaultUserAgent = `Mozilla/5.0 (Windows NT 10.0; Win64; x64)
	AppleWebKit/537.36 (KHTML, like Gecko)
	Chrome/124.0.0.0 Safari/537.36`
)

// defaultFlags contains the default Chromium command line flags,
// to be specific a [github.com/chromedp/chromedp.Flag] used by [Retina].
var defaultFlags = map[string]any{
	// disable-blink-features disables the AutomationControlled Blink feature,
	// which prevents pages from identifying automation through navigator.webdriver.
	"disable-blink-features": "AutomationControlled",

	// disable-dev-shm-usage prevents Chromium from using /dev/shm
	// for shared memory,
	// which can help avoid shared-memory limitations in containerized environments.
	"disable-dev-shm-usage": true,
}

// Option configures a [Retina] instance.
type Option func(*Retina)

// WithWindowSize configures the viewport dimensions
// used by [Retina].
func WithWindowSize(width, height int) Option {
	return func(r *Retina) {
		r.width = width
		r.height = height
	}
}

// WithTimeout configures [Retina]'s context timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(r *Retina) {
		r.timeout = timeout
	}
}

// WithUserAgent configures the [Retina] user agent.
func WithUserAgent(userAgent string) Option {
	return func(r *Retina) {
		r.userAgent = userAgent
	}
}

// WithFlags configures the Chrome flags used by [Retina].
func WithFlags(flags map[string]any) Option {
	return func(r *Retina) {
		for name, value := range flags {
			r.flags[name] = value
		}
	}
}

// New initializes a [Retina] instance with the default configurations,
// applies the provided [Option] values,
// and returns the configured instance.
func New(options ...Option) *Retina {
	r := &Retina{
		height:    defaultHeight,
		width:     defaultWidth,
		timeout:   defaultTimeout,
		userAgent: defaultUserAgent,
		flags:     maps.Clone(defaultFlags),
	}

	for _, option := range options {
		option(r)
	}

	return r
}
