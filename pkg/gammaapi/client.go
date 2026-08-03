package gamma

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"polymarket/internal/models"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	apiBaseURL  = "https://gamma-api.polymarket.com"
	isActive    = true
	isNotClosed = false
	cs2         = "Counter-strike-2"
	cs2TagID    = "100639"
	limit       = 100
)

type CSQueryParams struct {
	TagSlug    string    `url:"tag_slug,omitempty"`
	TagID      string    `url:"tag_id,omitempty"`
	Active     *bool     `url:"active,omitempty"`
	Closed     *bool     `url:"closed,omitempty"`
	EndDateMin time.Time `url:"end_date_min,omitempty"`
	Limit      int       `url:"limit,omitempty"`
	Offset     int       `url:"offset,omitempty"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		baseURL: apiBaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetEvents(ctx context.Context, params CSQueryParams) ([]models.Event, error) {
	values, err := query.Values(params)
	if err != nil {
		return nil, fmt.Errorf("pkg/client failed to build query: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/events?"+values.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("pkg/client failed to create http request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("pkg/client http request failed: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var eventDTOs []EventDTO
	if err := json.NewDecoder(response.Body).Decode(&eventDTOs); err != nil {
		return nil, fmt.Errorf("pkg/client failed to decode response body: %w", err)
	}

	events, err := mappedEvents(eventDTOs)
	if err != nil {
		return nil, fmt.Errorf("pkg/client failed to map events: %w", err)
	}

	return events, nil
}

func (c *Client) GetNCSEvents(ctx context.Context, N int) ([]models.Event, error) {
	active := isActive
	closed := isNotClosed
	params := CSQueryParams{
		TagSlug:    cs2,
		TagID:      cs2TagID,
		Active:     &active,
		Closed:     &closed,
		EndDateMin: time.Now().UTC(),
		Limit:      N,
	}

	return c.GetEvents(ctx, params)
}

func (c *Client) GetAllCSEvents(ctx context.Context) ([]models.Event, error) {
	active := isActive
	closed := isNotClosed

	var (
		offset    = 0
		allEvents []models.Event
	)

	for {
		params := CSQueryParams{
			TagSlug:    cs2,
			TagID:      cs2TagID,
			Active:     &active,
			Closed:     &closed,
			EndDateMin: time.Now().UTC(),
			Limit:      limit,
			Offset:     offset,
		}

		events, err := c.GetEvents(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("pkg/client failed fetching paginated events: %w", err)
		}

		allEvents = append(allEvents, events...)
		if len(events) < limit {
			break
		}

		offset += limit
	}

	return allEvents, nil
}
