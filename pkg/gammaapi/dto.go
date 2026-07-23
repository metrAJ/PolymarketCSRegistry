package gamma

type EventDTO struct {
	ID      string      `json:"id"`
	Title   string      `json:"title"`
	EndDate string      `json:"endDate"`
	Tags    []TagDTO    `json:"tags"`
	Markets []MarketDTO `json:"markets"`
}

type MarketDTO struct {
	ID              string  `json:"id"`
	Question        string  `json:"question"`
	EndDate         string  `json:"endDate"`
	GameStartTime   string  `json:"gameStartTime"`
	AcceptingOrders bool    `json:"acceptingOrders"`
	VolumeNum       float64 `json:"volumeNum"`
	LiquidityNum    float64 `json:"liquidityNum"`

	Outcomes      string `json:"outcomes"`
	OutcomePrices string `json:"outcomePrices"`
}

type TagDTO struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}
