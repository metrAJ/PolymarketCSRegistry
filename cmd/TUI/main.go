package main

import (
	"context"
	"fmt"
	"os"
	"polymarket/internal/client"
	"polymarket/internal/client/ui"
)

const urlStr = "http://localhost:3000/api/events"

func main() {
	events, err := client.FetchEvents(context.Background(), urlStr)
	if err != nil {
		fmt.Printf("Error fetching events: %v\n", err)
	}

	if len(events) == 0 {
		fmt.Println("No active events found.")
		return
	}

	appView := ui.NewView()
	appView.LoadEvents(events)

	if err := appView.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Terminal UI crashed: %v\n", err)
		os.Exit(1)
	}
}
