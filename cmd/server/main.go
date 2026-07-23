package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"polymarket/internal/config"
	gamma "polymarket/pkg/gammaapi"
	"time"

	"go.uber.org/zap"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	if cfg.Port == "" {
		cfg.Port = "3000"
	}
	// logger.Info("Starting server on port:" + cfg.Port)
	// log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))

	// gamma pkg test
	client := gamma.NewClient(logger)

	active := true
	closed := false
	params := gamma.CSQueryParams{
		TagSlug:    "Counter-strike-2",
		TagID:      "100639",
		Active:     &active,
		Closed:     &closed,
		EndDateMin: time.Now().UTC(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	events, err := client.GetEvents(ctx, params)

	if err != nil {
		logger.Error("Failed to get events", zap.Error(err))
		return
	}

	testJSON, _ := json.MarshalIndent(events[0], "", "  ")
	fmt.Println(string(testJSON))

}
