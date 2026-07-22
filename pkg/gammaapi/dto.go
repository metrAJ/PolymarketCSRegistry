package gamma

import "time"

type EventDTO struct {
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	EndDate time.Time   `json:"endDate"`
	Label   string      `json:"label"`
	Markets []MarketDTO `json:"markets"`
}

type MarketDTO struct {
	ID              string    `json:"id"`
	Question        string    `json:"question"`
	EndDate         time.Time `json:"endDate"`
	GameStartTime   time.Time `json:"gameStartTime"`
	AcceptingOrders bool      `json:"acceptingOrders"`
	VolumeNum       float64   `json:"volumeNum"`
	LiquidityNum    float64   `json:"liquidityNum"`

	Outcomes      string `json:"outcomes"`
	OutcomePrices string `json:"outcomePrices"`
}
