package main

import (
	"encoding/json"
	"os"
	"time"
)

type CacheData struct {
	loc_data Location      `json:"location"`
	weather  OpenMeteoResp `json:"weather"`
	time     time.Time     `json:"time"`
}

const cacheFile = "/tmp/weather_cache.json"

// init_cache erstellt die Cache-Datei, falls sie noch nicht existiert.
// Falls sie existiert, macht die Funktion nichts.
func init_cache() (bool, error) {
	// Prüfen, ob Datei schon existiert
	if _, err := os.Stat(cacheFile); err == nil {
		// existiert schon → nichts tun
		return true, nil
	}

	// Leeren Cache anlegen
	empty := CacheData{}

	file, err := os.Create(cacheFile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(empty); err != nil {
		return false, err
	}

	return false, nil
}

func update_cache(loc_data Location, weather_data OpenMeteoResp) {
	// Datei anlegen/überschreiben
	file, err := os.Create(cacheFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Daten in Struct packen
	data := CacheData{
		loc_data: loc_data,
		weather:  weather_data,
		time:     time.Now(),
	}

	// JSON-Encoder vorbereiten
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // schön formatiert

	// schreiben
	if err := encoder.Encode(data); err != nil {
		panic(err)
	}

}

func load_cache() (*CacheData, error) {
	file, err := os.Open(cacheFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data CacheData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func time_of_cache(cache *CacheData) time.Duration {
	if cache == nil {
		return 0
	}
	return time.Since(cache.time)
}
