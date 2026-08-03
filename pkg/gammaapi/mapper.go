package gamma

// Helper to map EventDTOs to domain Events

import (
	"encoding/json"
	"fmt"
	"polymarket/internal/models"
	"strconv"
	"time"
)

func mappedEvents(eventDTOs []EventDTO) ([]models.Event, error) {
	events := make([]models.Event, 0, len(eventDTOs))

	for _, dto := range eventDTOs {
		domainEvent, err := dto.ToDomainModel()
		if err != nil {
			return nil, fmt.Errorf("failed to map EventDTO to domain Event: %w", err)
		}

		events = append(events, domainEvent)
	}

	return events, nil
}

func (dto *EventDTO) ToDomainModel() (models.Event, error) {
	endDate, err := parseTime(dto.EndDate)
	if err != nil {
		return models.Event{}, fmt.Errorf("failed to parse end date for event %s: %w", dto.ID, err)
	}

	modelEvent := models.Event{
		ID:      dto.ID,
		Title:   dto.Title,
		EndDate: endDate,
		Tags:    parseTags(dto.Tags),
		Markets: make([]models.Market, 0, len(dto.Markets)),
	}

	for _, m := range dto.Markets {
		outcomes, err := parseOutcomes(m.Outcomes, m.OutcomePrices)
		if err != nil {
			return models.Event{}, fmt.Errorf("failed to parse outcomes for market %s: %w", m.ID, err)
		}

		endDate, err := parseTime(m.EndDate)
		if err != nil {
			return models.Event{}, fmt.Errorf("failed to parse end date for market %s: %w", m.ID, err)
		}

		gameStartTime, err := parseTime(m.GameStartTime)
		if err != nil {
			return models.Event{}, fmt.Errorf("failed to parse game start time for market %s: %w", m.ID, err)
		}

		modelEvent.Markets = append(modelEvent.Markets, models.Market{
			ID:              m.ID,
			Question:        m.Question,
			EndDate:         endDate,
			GameStartTime:   gameStartTime,
			AcceptingOrders: m.AcceptingOrders,
			VolumeNum:       m.VolumeNum,
			LiquidityNum:    m.LiquidityNum,
			Outcomes:        outcomes,
		})
	}

	return modelEvent, nil
}

func parseOutcomes(outcomesJSON, outcomePricesJSON string) ([]models.Outcome, error) {
	var (
		names         []string
		outcomePrices []string
	)

	err := json.Unmarshal([]byte(outcomesJSON), &names)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal outcomes JSON: %w", err)
	}

	if outcomePricesJSON != "" && outcomePricesJSON != "null" {
		err := json.Unmarshal([]byte(outcomePricesJSON), &outcomePrices)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal outcome prices JSON: %w", err)
		}
	}

	if len(names) > 0 && len(outcomePrices) == 0 {
		outcomePrices = make([]string, len(names))
		for i := range outcomePrices {
			outcomePrices[i] = "0"
		}
	}

	if len(names) != len(outcomePrices) {
		return nil, fmt.Errorf("mismatched lengths: %d outcomes and %d outcome prices", len(names), len(outcomePrices))
	}

	outcomes := make([]models.Outcome, 0, len(names))

	for i, name := range names {
		price, err := strconv.ParseFloat(outcomePrices[i], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse outcome price '%s': %w", outcomePrices[i], err)
		}

		outcomes = append(outcomes, models.Outcome{
			Outcome:       name,
			OutcomePrices: price,
		})
	}

	return outcomes, nil
}

func parseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("empty time string")
	}

	layout := "2006-01-02 15:04:05-07"

	t, err := time.Parse(layout, timeStr)
	if err != nil {
		t, fallbackErr := time.Parse(time.RFC3339, timeStr)
		if fallbackErr != nil {
			return time.Time{}, fmt.Errorf("failed to parse time string '%s': %w", timeStr, err)
		}

		return t, nil
	}

	return t, nil
}

func parseTags(tagDTOs []TagDTO) []string {
	tags := make([]string, 0, len(tagDTOs))

	for _, tag := range tagDTOs {
		if tag.Label != "" {
			tags = append(tags, tag.Label)
		}
	}

	return tags
}
