package ui

import (
	"fmt"
	"polymarket/internal/client"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type View struct {
	app     *tview.Application
	layout  *tview.Flex
	list    *tview.List
	details *tview.TextView
	events  []client.Event
}

func NewView() *View {
	v := &View{
		app:     tview.NewApplication(),
		layout:  tview.NewFlex(),
		list:    tview.NewList().ShowSecondaryText(false),
		details: tview.NewTextView().SetDynamicColors(true).SetWrap(true),
	}

	v.setupComponents()
	v.setupKeybindings()

	return v
}

func (v *View) setupComponents() { // Split screen on 2 parts
	v.list.SetBorder(true).SetTitle(" Active Events ")
	v.details.SetBorder(true).SetTitle(" Event Details ")

	v.layout.AddItem(v.list, 0, 1, true).AddItem(v.details, 0, 2, false) // 1 to 2 proportion
}

func (v *View) setupKeybindings() { // Key capture
	v.app.SetInputCapture(func(tcellEvent *tcell.EventKey) *tcell.EventKey {
		switch tcellEvent.Key() {
		case tcell.KeyTab: // Change focus between parts
			if v.list.HasFocus() {
				v.app.SetFocus(v.details)
			} else {
				v.app.SetFocus(v.list)
			}

			return nil

		case tcell.KeyCtrlC, tcell.KeyEscape:
			v.app.Stop()
			return nil
		}

		return tcellEvent
	})
	v.list.SetChangedFunc(func(index int, _, _ string, _ rune) {
		v.renderDetails(v.events[index])
	})
}

func (v *View) LoadEvents(events []client.Event) {
	v.events = events
	v.list.Clear()

	for _, event := range events {
		v.list.AddItem(event.Title, "", 0, nil)
	}

	if len(events) > 0 {
		v.renderDetails(events[0])
	}
}

func (v *View) renderDetails(event client.Event) {
	var sb strings.Builder

	endDate := event.EndDate
	if t, err := time.Parse(time.RFC3339, event.EndDate); err == nil {
		endDate = t.Format("Jan 02, 2006 15:04 UTC")
	}

	fmt.Fprintf(&sb, "[yellow]TITLE:[-] %s\n", event.Title)
	fmt.Fprintf(&sb, "[yellow]CLOSES:[-] %s\n", endDate)
	fmt.Fprintf(&sb, "[yellow]TAGS:[-] %s\n\n", strings.Join(event.Tags, ", "))

	for i, market := range event.Markets {
		fmt.Fprintf(&sb, "[green]--- MARKET %d ---[-]\n", i+1)
		fmt.Fprintf(&sb, "QUESTION:  %s\n", market.Question)
		fmt.Fprintf(&sb, "VOLUME:    $%.2f\n", market.VolumeNum)
		fmt.Fprintf(&sb, "LIQUIDITY: $%.2f\n\n", market.LiquidityNum)
		sb.WriteString("[blue]OUTCOMES:[-]\n")

		for _, outcome := range market.Outcomes {
			fmt.Fprintf(&sb, "  - %s: %g\n", outcome.Outcome, outcome.OutcomePrices)
		}

		sb.WriteString("\n")
	}

	v.details.SetText(sb.String())
}

func (v *View) Start() error {
	return v.app.SetRoot(v.layout, true).EnableMouse(true).Run()
}

func (v *View) Stop() {
	v.app.Stop()
}
