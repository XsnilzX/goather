package main

import (
	"testing"
	"time"
)

func TestConvertToHourlyParsesOpenMeteoTimes(t *testing.T) {
	t.Parallel()

	var api OpenMeteoResp
	api.Hourly.Time = []string{"2026-04-24T16:00", "2026-04-24T17:00"}
	api.Hourly.Temperature2m = []float64{15.5, 16.2}
	api.Hourly.WeatherCode = []int{1, 2}

	hours := convertToHourly(api)
	if len(hours) != 2 {
		t.Fatalf("expected 2 hourly entries, got %d", len(hours))
	}
	if got := hours[0].Time.Format("2006-01-02T15:04"); got != "2026-04-24T16:00" {
		t.Fatalf("unexpected parsed time %q", got)
	}
}

func TestParseForecastTimeAcceptsRFC3339(t *testing.T) {
	t.Parallel()

	parsed, err := parseForecastTime("2026-04-24T16:00:00Z")
	if err != nil {
		t.Fatalf("expected RFC3339 fallback to parse: %v", err)
	}
	if !parsed.Equal(time.Date(2026, 4, 24, 16, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected parsed RFC3339 time: %s", parsed)
	}
}
