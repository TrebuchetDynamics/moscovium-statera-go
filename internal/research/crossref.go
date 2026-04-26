package research

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultCrossrefBaseURL = "https://api.crossref.org"

type Work struct {
	DOI     string
	Title   string
	Journal string
	Year    int
}

type CrossrefClient struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

type CrossrefOption func(*CrossrefClient)

func NewCrossrefClient(options ...CrossrefOption) *CrossrefClient {
	client := &CrossrefClient{
		baseURL: defaultCrossrefBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: "moscovium-statera-go/0.1 (+https://github.com/TrebuchetDynamics/moscovium-statera-go)",
	}
	for _, option := range options {
		option(client)
	}
	return client
}

func WithBaseURL(baseURL string) CrossrefOption {
	return func(client *CrossrefClient) {
		client.baseURL = strings.TrimRight(baseURL, "/")
	}
}

func WithUserAgent(userAgent string) CrossrefOption {
	return func(client *CrossrefClient) {
		client.userAgent = userAgent
	}
}

func WithHTTPClient(httpClient *http.Client) CrossrefOption {
	return func(client *CrossrefClient) {
		client.httpClient = httpClient
	}
}

func (c *CrossrefClient) FetchWork(ctx context.Context, doi string) (Work, error) {
	if strings.TrimSpace(doi) == "" {
		return Work{}, fmt.Errorf("doi is required")
	}

	url := c.baseURL + "/works/" + strings.TrimSpace(doi)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Work{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Work{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Work{}, fmt.Errorf("crossref returned HTTP %d", resp.StatusCode)
	}

	var payload struct {
		Status  string `json:"status"`
		Message struct {
			DOI            string     `json:"DOI"`
			Title          []string   `json:"title"`
			ContainerTitle []string   `json:"container-title"`
			Issued         dateIssued `json:"issued"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return Work{}, err
	}
	if payload.Status != "" && payload.Status != "ok" {
		return Work{}, fmt.Errorf("crossref status %q", payload.Status)
	}

	return Work{
		DOI:     payload.Message.DOI,
		Title:   first(payload.Message.Title),
		Journal: first(payload.Message.ContainerTitle),
		Year:    payload.Message.Issued.Year(),
	}, nil
}

type dateIssued struct {
	DateParts [][]int `json:"date-parts"`
}

func (d dateIssued) Year() int {
	if len(d.DateParts) == 0 || len(d.DateParts[0]) == 0 {
		return 0
	}
	return d.DateParts[0][0]
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
