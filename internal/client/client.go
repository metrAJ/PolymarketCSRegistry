package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Event struct {
	Title   string   `json:"title"`
	EndDate string   `json:"endDate"`
	Tags    []string `json:"tags"`
	Markets []Market `json:"markets"`
}

type Market struct {
	Question     string    `json:"question"`
	VolumeNum    float64   `json:"volumeNum"`
	LiquidityNum float64   `json:"liquidityNum"`
	Outcomes     []Outcome `json:"outcomes"`
}

type Outcome struct {
	Outcome       string  `json:"outcome"`
	OutcomePrices float64 `json:"outcomePrices"`
}

func FetchEvents(url string) ([]Event, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("json.Decode: %w", err)
	}

	return events, nil
}
