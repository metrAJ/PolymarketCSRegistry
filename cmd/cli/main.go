package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"polymarket/internal/client"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"
)

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func main() {
	urlStr := "http://localhost:3000/api/events"

	events, err := client.FetchEvents(urlStr)
	if err != nil {
		fmt.Printf("Error fetching events: %v\n", err)
	}

	if len(events) == 0 {
		fmt.Println("No active events found.")
		return
	}

	// I did not really imagine on how to make a goodlloking output, but enhanced with ai.

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("\n FOUND %d ACTIVE COUNTER-STRIKE EVENTS \n", len(events))
	fmt.Println(strings.Repeat("=", 60))

	for i, event := range events {
		endDate := event.EndDate
		if t, err := time.Parse(time.RFC3339, event.EndDate); err == nil {
			endDate = t.Format("Jan 02, 2006 15:04 UTC")
		}

		fmt.Fprintf(w, "EVENT (%d/%d):\t%s\n", i+1, len(events), event.Title)
		fmt.Fprintf(w, "CLOSES:\t%s\n", endDate)
		fmt.Fprintf(w, "TAGS:\t%s\n", strings.Join(event.Tags, ", "))
		fmt.Fprintln(w, "\t")

		for j, market := range event.Markets {
			fmt.Fprintf(w, "  -> MARKET %d:\t%s\n", j+1, market.Question)
			fmt.Fprintf(w, "     VOLUME:\t$%.2f\n", market.VolumeNum)
			fmt.Fprintf(w, "     LIQUIDITY:\t$%.2f\n", market.LiquidityNum)
			fmt.Fprintln(w, "     OUTCOMES:\t")

			for _, outcome := range market.Outcomes {
				fmt.Fprintf(w, "\t- %s:\t%g\n", outcome.Outcome, outcome.OutcomePrices)
			}

			fmt.Fprintln(w, "\t")
		}

		w.Flush()
		fmt.Println(strings.Repeat("-", 60))

		if i < len(events)-1 {
			fmt.Print("Press [ENTER] to load the next event, or [CTRL+C] to exit...")
			scanner.Scan()

			clearScreen()
		}
	}

	fmt.Println("\nAll events loaded.")
}
