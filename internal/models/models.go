package models

import "time"

// Raw models without JSON tags

// Event struct, will contain all necessary for us info about events
type Event struct {
	ID      string
	Title   string
	EndDate time.Time
	Tags    []string
	Markets []Market
}

// Market struct, will contain all necessary for us info about markets
type Market struct {
	ID              string
	Question        string
	EndDate         time.Time
	Outcomes        []Outcome
	GameStartTime   time.Time
	AcceptingOrders bool
	VolumeNum       float64
	LiquidityNum    float64
}

// I thought about moving outcomes and their chances to different struct in case there will be more outcomes and more outcome types, like in multi markets with tournament bets. Plus we might need new fileds if we want to iplement trading functional in future.
type Outcome struct {
	Outcome       string
	OutcomePrices float64
}
