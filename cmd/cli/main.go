package main

import (
	"context"
	"fmt"
	"os"
	"polymarket/internal/client"
	"polymarket/internal/client/cli"
)

const urlStr = "http://localhost:3000/api/events"

func main() {
	events, err := client.FetchEvents(context.Background(), urlStr)
	if err != nil {
		fmt.Printf("Error fetching events: %v\n", err)
		os.Exit(1)
	}

	if len(events) == 0 {
		fmt.Println("No active events found.")
		return
	}

	cliPrinter := cli.NewPrinter(os.Stdout, os.Stdin)
	cliPrinter.Render(events)
	fmt.Println("\nAll events loaded.")
}
