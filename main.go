package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

const HOURS int8 = 8

type QuickshellOut struct {
	Display string `json:"display"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
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
		_ = json.NewEncoder(os.Stdout).Encode(QuickshellOut{
			Display: "❌ Error",
			Tooltip: "Location error: " + err.Error(),
			Class:   "error",
		})
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
	// make cache if not already exists
	cacheExists, err := init_cache()
	if err != nil {
		panic(err)
	}

	var cache *CacheData
	if cacheExists {
		cache2, err := load_cache()
		if err != nil {
			_ = json.NewEncoder(os.Stdout).Encode(QuickshellOut{
				Display: "❌ Error",
				Tooltip: "Cache error: " + err.Error(),
				Class:   "error",
			})
			return
		}
		cache = cache2
	}

	var api *OpenMeteoResp
	cacheFresh := cache != nil && time.Since(cache.time) <= time.Hour
	if cacheFresh {
		api = &cache.weather
	} else {
		api2, err := FetchWeather(ctx, lat, lon, 6) // z. B. 6 Stunden Vorhersage
		if err != nil {
			_ = json.NewEncoder(os.Stdout).Encode(QuickshellOut{
				Display: "❌ Error",
				Tooltip: "Weather error: " + err.Error(),
				Class:   "error",
			})
			return
		}
		api = api2
		update_cache(loc, *api)
	}

	// === Text (kurz) für Quickshell ===
	text := fmt.Sprintf("%s %.0f°C", iconFor(api.Current.WeatherCode), api.Current.Temperature2m)
	locationLine := fmt.Sprintf("%s, %s", loc.City, loc.Country)
	if loc.Region != "" && loc.Region != loc.City {
		locationLine = fmt.Sprintf("%s, %s, %s", loc.City, loc.Region, loc.Country)
	}

	tooltipLines := []string{
		locationLine,
		descriptionFor(api.Current.WeatherCode),
		fmt.Sprintf("🌡️ Temperature: %.0f°C (feels %.0f°C)", api.Current.Temperature2m, api.Current.ApparentTemperature),
		fmt.Sprintf("💧 Humidity: %d%%", int(api.Current.RelativeHumidity2)),
		fmt.Sprintf("💨 Wind: %.0f km/h", api.Current.WindSpeed10m),
	}
	if forecast := formatHourlyForecast(convertToHourly(*api), 6); forecast != "" {
		tooltipLines = append(tooltipLines, forecast)
	}
	tooltipLines = append(tooltipLines, fmt.Sprintf("Updated: %s", time.Now().Format("15:04")))
	tooltip := strings.Join(tooltipLines, "\n")

	_ = json.NewEncoder(os.Stdout).Encode(QuickshellOut{
		Display: text,
		Tooltip: tooltip,
		Class:   classFor(api.Current.WeatherCode),
	})
}
