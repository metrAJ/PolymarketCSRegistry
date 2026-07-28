package gamma

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"polymarket/internal/models"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	apiBaseURL = "https://gamma-api.polymarket.com"
)

type CSQueryParams struct {
	TagSlug    string    `url:"tag_slug,omitempty"`
	TagID      string    `url:"tag_id,omitempty"`
	Active     *bool     `url:"active,omitempty"`
	Closed     *bool     `url:"closed,omitempty"`
	EndDateMin time.Time `url:"end_date_min,omitempty"`
}

const (
	isActive    = true
	isNotClosed = false
	cs2         = "Counter-strike-2"
	cs2TagID    = "100639"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
}

func NewClient(logger *slog.Logger) *Client {
	return &Client{
		baseURL:    apiBaseURL,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

func (c *Client) GetEvents(ctx context.Context, params CSQueryParams) ([]models.Event, error) {
	values, err := query.Values(params)
	if err != nil {
		c.logger.Error("pkg/client failed to build querry", "error", err)
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/events?"+values.Encode(), nil)
	if err != nil {
		c.logger.Error("pkg/client failed to create http request", "error", err)
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.logger.Error("pkg/client http request failed", "error", err)
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected status code: %d", response.StatusCode)

		c.logger.Error("gamma api returned non-200 status", "status", response.StatusCode)

		return nil, err
	}

	var eventDTOs []EventDTO
	if err := json.NewDecoder(response.Body).Decode(&eventDTOs); err != nil {
		c.logger.Error("failed to decode response body", "error", err)
		return nil, err
	}

	events, err := mappedEvents(eventDTOs)
	if err != nil {
		c.logger.Error("failed to map event DTOs", "error", err)
		return nil, fmt.Errorf("failed to map events: %w", err)
	}

	return events, nil
}

func (c *Client) GetCSEvents(ctx context.Context) ([]models.Event, error) {
	active := isActive
	closed := isNotClosed
	params := CSQueryParams{
		TagSlug:    cs2,
		TagID:      cs2TagID,
		Active:     &active,
		Closed:     &closed,
		EndDateMin: time.Now().UTC(),
	}

	return c.GetEvents(ctx, params)
}
