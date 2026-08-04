package gamma

import (
	"polymarket/internal/models"
	"reflect"
	"testing"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name    string
		timeStr string
		wantErr bool
	}{
		{
			name:    "valid format",
			timeStr: "2026-08-03 15:04:05+03",
			wantErr: false,
		},
		{
			name:    "empty string",
			timeStr: "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			timeStr: "123123=123-123",
			wantErr: true,
		},
		{
			name:    "valid RFC3339 format",
			timeStr: "2026-08-03T15:04:05Z",
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseTime(test.timeStr)
			if (err != nil) != test.wantErr {
				t.Errorf("parseTime() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestParseOutcomes(t *testing.T) {
	tests := []struct {
		name       string
		namesJSON  string
		pricesJSON string
		want       []models.Outcome
		wantErr    bool
	}{
		{
			name:       "valid outcomes and prices",
			namesJSON:  `["Yes", "No"]`,
			pricesJSON: `["0.75", "0.25"]`,
			want: []models.Outcome{
				{Outcome: "Yes", OutcomePrices: 0.75},
				{Outcome: "No", OutcomePrices: 0.25},
			},
			wantErr: false,
		},
		{
			name:       "missing prices, went to 0",
			namesJSON:  `["Navi", "FaZe"]`,
			pricesJSON: "",
			want: []models.Outcome{
				{Outcome: "Navi", OutcomePrices: 0.0},
				{Outcome: "FaZe", OutcomePrices: 0.0},
			},
			wantErr: false,
		},
		{
			name:       "mismatched lengths",
			namesJSON:  `["Yes", "No", "Maybe"]`,
			pricesJSON: `["0.5", "0.5"]`,
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "invalid float parsing",
			namesJSON:  `["Yes"]`,
			pricesJSON: `["300-bucks"]`,
			want:       nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOutcomes(tt.namesJSON, tt.pricesJSON)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOutcomes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOutcomes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseTags(t *testing.T) {
	tests := []struct {
		name string
		tags []TagDTO
		want []string
	}{
		{
			name: "valid labels",
			tags: []TagDTO{
				{Label: "Sports"},
				{Label: "Esports"},
			},
			want: []string{"Sports", "Esports"},
		},
		{
			name: "empty labels",
			tags: []TagDTO{
				{Label: "Crypto"},
				{Label: ""},
				{Label: "Politics"},
			},
			want: []string{"Crypto", "Politics"},
		},
		{
			name: "empty input",
			tags: []TagDTO{},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseTags(tt.tags)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTags() got = %v, want %v", got, tt.want)
			}
		})
	}
}
