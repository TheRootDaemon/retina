package retina

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

// Retina represents the main entry point
// for managing browser automation.
// It runs as a headless browser by default
// and might be configured to suit different use cases.
type Retina struct {
	height int
	width  int

	timeout time.Duration

	userAgent string
	flags     map[string]any

	ctx    context.Context
	cancel context.CancelFunc
}

// SetupBrowser creates and configures a new browser context
// using the configuration of [Retina].
//
// The returned context is ready for running browser automation tasks.
//
// The returned cancel function must be called
// to release the resources associated with the browser context.
func (r *Retina) SetupBrowser(parent context.Context) (context.Context, context.CancelFunc) {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],

		chromedp.DisableGPU,
		chromedp.Headless,
		chromedp.NoSandbox,

		chromedp.UserAgent(r.userAgent),
		chromedp.WindowSize(r.width, r.height),
	)

	for name, value := range r.flags {
		opts = append(opts, chromedp.Flag(name, value))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(parent, opts...)
	ctx, timeoutCancel := context.WithTimeout(allocCtx, r.timeout)
	taskCtx, taskCancel := chromedp.NewContext(ctx)

	cancel := func() {
		taskCancel()
		timeoutCancel()
		allocCancel()
	}

	return taskCtx, cancel
}
