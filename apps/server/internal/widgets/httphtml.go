package widgets

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/network"
)

// HTTPHTMLWidget scrapes HTML and extracts text
type HTTPHTMLWidget struct{}

func (w *HTTPHTMLWidget) Type() string {
	return "httpHtml"
}

func (w *HTTPHTMLWidget) CacheTTL() time.Duration {
	return 120 * time.Second
}

func (w *HTTPHTMLWidget) Execute(ctx context.Context, cfg config.Widget, client *network.Client) (*Result, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", cfg.URL, nil)
	if err != nil {
		return &Result{Error: err.Error(), State: "error"}, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return &Result{Error: err.Error(), State: "error"}, nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Result{Error: err.Error(), State: "error"}, nil
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return &Result{Error: "failed to parse HTML: " + err.Error(), State: "error"}, nil
	}

	var value string
	if cfg.Attribute != "" {
		value, _ = doc.Find(cfg.Selector).Attr(cfg.Attribute)
	} else {
		value = strings.TrimSpace(doc.Find(cfg.Selector).First().Text())
	}

	if value == "" {
		return &Result{
			Error:      "no content found for selector: " + cfg.Selector,
			State:      "error",
			LastUpdate: time.Now(),
		}, nil
	}

	state := "good"
	if strings.Contains(strings.ToLower(value), "error") ||
		strings.Contains(strings.ToLower(value), "down") {
		state = "error"
	}

	return &Result{
		Label:      cfg.Label,
		Value:      value,
		Formatted:  value,
		State:      state,
		LastUpdate: time.Now(),
	}, nil
}
