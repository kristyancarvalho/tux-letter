package fetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxBodyBytes = 8 << 20

type Response struct {
	URL         string
	Status      int
	ContentType string
	Body        []byte
}

type Fetcher interface {
	Get(ctx context.Context, url string) (*Response, error)
}

type Client struct {
	http      *http.Client
	userAgent string
}

func New(timeout time.Duration, userAgent string) *Client {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	if userAgent == "" {
		userAgent = "tux-letter/3.0"
	}
	return &Client{
		http:      &http.Client{Timeout: timeout},
		userAgent: userAgent,
	}
}

func (c *Client) Get(ctx context.Context, url string) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml,application/rss+xml,application/atom+xml;q=0.9,*/*;q=0.8")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch %s: status %d", url, resp.StatusCode)
	}

	return &Response{
		URL:         resp.Request.URL.String(),
		Status:      resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Body:        body,
	}, nil
}
