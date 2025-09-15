package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const HOURS int8 = 8

type WaybarOut struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// === Location bestimmen ===

	loc, err := GetLocation(
		ctx,
		WithOverallTimeout(1800*time.Millisecond),
		WithPerReqTimeout(1000*time.Millisecond),
		WithPreferFastest(false),
	)
	if err != nil {
		_ = json.NewEncoder(os.Stdout).Encode(WaybarOut{Text: "⚠️", Tooltip: "Location error: " + err.Error()})
		return
	}
	lat, lon := loc.Lat, loc.Lon

	// Demo: Falls du noch keine Location-Funktion einbindest, nimm vorerst Berlin:
	/* lat, lon := 52.52, 13.41
	loc := struct {
		City, Region, Country, Source string
	}{City: "Berlin", Region: "Berlin", Country: "Deutschland", Source: "demo"} */

	// === Wetter holen ===
	// implement weather cache
	
	cache_time := get_time_of_cache()
	if cache_time = nil {
		// create new Cache
	}
	
	curent := time.Now()
	 

	
	api, err := FetchWeather(ctx, lat, lon, 6) // z. B. 6 Stunden Vorhersage
	if err != nil {
		_ = json.NewEncoder(os.Stdout).Encode(WaybarOut{Text: "⚠️", Tooltip: "Weather error: " + err.Error()})
		return
	}

	// === Text (kurz) für Waybar ===
	text := fmt.Sprintf("%s %.0f°C", iconFor(api.Current.WeatherCode), api.Current.Temperature2m)

	// === Tooltip wie früher, nur in Go ===
	tooltip := fmt.Sprintf(
		"📍 %s, %s, %s\n"+
			"<span size='xx-large'>%.0f°C</span>\n"+
			"<big>%s %s</big>\n"+
			"Gefühlte: %.0f°C\n"+
			"Feuchtigkeit: %d%%\n"+
			"Wind: %.0f km/h%s",
		loc.City, loc.Region, loc.Country,
		api.Current.Temperature2m,
		iconFor(api.Current.WeatherCode),
		descriptionFor(api.Current.WeatherCode),
		api.Current.ApparentTemperature,
		int(api.Current.RelativeHumidity2),
		api.Current.WindSpeed10m,
		formatHourlyForecast(convertToHourly(*api), 5), // z. B. 5 Stunden anhängen
	)

	_ = json.NewEncoder(os.Stdout).Encode(WaybarOut{
		Text:    text,
		Tooltip: tooltip,
	})
}
