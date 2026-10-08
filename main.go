package main

import (
	"context"
	"fmt"

	"github.com/therootdaemon/retina/internal/retina"
)

func main() {
	b := retina.New()
	b.SetupBrowser(context.Background())
	fmt.Println(b.Search(&retina.DuckDuckGo{}, "javascript"))
}
