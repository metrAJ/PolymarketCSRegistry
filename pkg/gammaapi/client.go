package gamma

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"polymarket/internal/models"
	"time"

	"github.com/google/go-querystring/query"
	"go.uber.org/zap"
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

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

func NewClient(logger *zap.Logger) *Client {
	return &Client{
		baseURL:    apiBaseURL,
		httpClient: &http.Client{},
		logger:     logger,
	}
}

func (c *Client) GetEvents(ctx context.Context, params CSQueryParams) ([]models.Event, error) {
	values, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/events?"+values.Encode(), nil)

	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var eventDTOs []EventDTO
	if err := json.NewDecoder(response.Body).Decode(&eventDTOs); err != nil {
		return nil, err
	}

	events, err := mappedEvents(eventDTOs)
	if err != nil {
		return nil, fmt.Errorf("failed to map events: %w", err)
	}
	return events, nil
}
