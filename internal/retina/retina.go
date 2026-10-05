package retina

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/therootdaemon/retina/logger"
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
// The browser context is stored internally
// and can be reused by subsequent operations.
//
// [Retina.SetupBrowser] must be called before using browser operations
// such as [Retina.Search], otherwise the operations might end up using zero values.
// Call [Retina.Close] to cleanup resources.
func (r *Retina) SetupBrowser(parent context.Context) {
	logger.Info(
		"retina: setting up browser: dimensions=%dx%d timeout=%v user_agent=%s flags=%v",
		r.width,
		r.height,
		r.timeout,
		r.userAgent,
		r.flags,
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(parent, r.allocateOptions()...)
	ctx, timeoutCancel := context.WithTimeout(allocCtx, r.timeout)
	taskCtx, taskCancel := chromedp.NewContext(ctx)

	cancel := func() {
		taskCancel()
		timeoutCancel()
		allocCancel()
	}

	r.ctx = taskCtx
	r.cancel = cancel

	logger.Info("retina: setup complete")
}

// Close releases the resources associated with the browser.
func (r *Retina) Close() {
	if r.cancel == nil {
		logger.Debug("retina: browser is not running")
		return
	}

	logger.Info("retina: closing browser")
	r.cancel()
	logger.Info("retina: browser closed")
}

// allocateOptions returns the Chrome execution options
// configured for [Retina].
func (r *Retina) allocateOptions() []chromedp.ExecAllocatorOption {
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

	return opts
}
