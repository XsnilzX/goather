package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OpenMeteoResp ist die Antwortstruktur der API
type OpenMeteoResp struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Time                string  `json:"time"`
		Temperature2m       float64 `json:"temperature_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		WindSpeed10m        float64 `json:"wind_speed_10m"`
		RelativeHumidity2   float64 `json:"relative_humidity_2m"`
		WeatherCode         int     `json:"weather_code"`
	} `json:"current"`
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
		WeatherCode   []int     `json:"weather_code"`
	} `json:"hourly"`
}

// FetchWeather ruft die Open-Meteo API ab
func FetchWeather(ctx context.Context, lat, lon float64, hours int8) (*OpenMeteoResp, error) {
	if hours < 1 {
		hours = 1
	}
	if hours > 24 {
		hours = 24
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f"+
			"&current=temperature_2m,apparent_temperature,wind_speed_10m,relative_humidity_2m,weather_code"+
			"&hourly=temperature_2m,weather_code&forecast_hours=%d&timezone=auto",
		lat, lon, hours,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "waybar-openmeteo-go/1.0")

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var api OpenMeteoResp
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		return nil, err
	}
	return &api, nil
}
