// discord.go
package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TheRootDaemon/watchman/detector"
)

func SendReport(hook string, report detector.Report) error {
	content, err := json.MarshalIndent(report, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	payload := struct {
		Content string `json:"content"`
	}{
		Content: fmt.Sprintf("**Down Alert!**\n```json\n%s```", string(content)),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	r, err := http.Post(hook, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer r.Body.Close()
	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return fmt.Errorf("discord returned unexpected status: %w", err)
	}

	return nil
}
