package main

import (
	"fmt"
	"strings"
	"time"
)

// Icons für Open-Meteo Wettercodes
var weatherIcons = map[int]string{
	0: "☀️", 1: "🌤️", 2: "⛅", 3: "☁️",
	45: "🌫️", 48: "🌫️",
	51: "🌦️", 53: "🌦️", 55: "🌧️",
	61: "🌦️", 63: "🌧️", 65: "🌧️",
	71: "🌨️", 73: "🌨️", 75: "❄️", 77: "🌨️",
	80: "🌧️", 81: "🌧️", 82: "🌧️",
	95: "⛈️", 96: "⛈️", 99: "⛈️",
}

// Deutsche Beschreibungen
var weatherDescriptions = map[int]string{
	0: "Klarer Himmel", 1: "Überwiegend klar", 2: "Teilweise bewölkt", 3: "Bedeckt",
	45: "Nebel", 48: "Reif-Nebel",
	51: "Leichter Nieselregen", 53: "Mäßiger Nieselregen", 55: "Starker Nieselregen",
	61: "Leichter Regen", 63: "Mäßiger Regen", 65: "Starker Regen",
	71: "Leichter Schneefall", 73: "Mäßiger Schneefall", 75: "Starker Schneefall",
	80: "Regenschauer", 81: "Heftiger Regenschauer", 82: "Starker Regenschauer",
	95: "Gewitter", 96: "Gewitter mit Hagel", 99: "Schweres Gewitter mit Hagel",
}

func iconFor(code int) string {
	if s, ok := weatherIcons[code]; ok {
		return s
	}
	return "❓"
}

func descriptionFor(code int) string {
	if s, ok := weatherDescriptions[code]; ok {
		return s
	}
	return "Unbekannt"
}

func classFor(code int) string {
	switch code {
	case 0, 1:
		return "clear"
	case 2, 3:
		return "cloudy"
	case 45, 48:
		return "fog"
	case 51, 53, 55:
		return "drizzle"
	case 61, 63, 65, 80, 81, 82:
		return "rain"
	case 71, 73, 75:
		return "snow"
	case 95, 96, 99:
		return "storm"
	default:
		return "cloudy"
	}
}

// Kompakter Typ nur für das Tooltip
type HourlyForecast struct {
	Time        time.Time
	Temperature float64
	Code        int
}

// Konvertiert OpenMeteoResp -> []HourlyForecast (robust gegen unterschiedliche Längen)
func convertToHourly(api OpenMeteoResp) []HourlyForecast {
	n := len(api.Hourly.Time)
	if len(api.Hourly.Temperature2m) < n {
		n = len(api.Hourly.Temperature2m)
	}
	if len(api.Hourly.WeatherCode) < n {
		n = len(api.Hourly.WeatherCode)
	}
	out := make([]HourlyForecast, 0, n)
	for i := 0; i < n; i++ {
		t, err := time.Parse(time.RFC3339, api.Hourly.Time[i])
		if err != nil {
			continue
		}
		out = append(out, HourlyForecast{
			Time:        t,
			Temperature: api.Hourly.Temperature2m[i],
			Code:        api.Hourly.WeatherCode[i],
		})
	}
	return out
}

// Baut den Mehrstunden-Block (z. B. 5 Stunden) fürs Tooltip
func formatHourlyForecast(hours []HourlyForecast, count int) string {
	if len(hours) == 0 {
		return ""
	}
	if count > len(hours) {
		count = len(hours)
	}
	var b strings.Builder
	b.WriteString("⏰ Next 6 hours:")
	for i := 0; i < count; i++ {
		h := hours[i]
		b.WriteString(fmt.Sprintf("\n%s %s %.0f°C", h.Time.Format("15:04"), iconFor(h.Code), h.Temperature))
	}
	return b.String()
}
