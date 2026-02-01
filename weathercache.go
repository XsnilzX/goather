package main

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"time"
)

type CacheData struct {
	SchemaVersion int           `json:"schema_version"`
	CachedAt      time.Time     `json:"cached_at"`
	ExpiresAt     time.Time     `json:"expires_at"`
	Lat           float64       `json:"lat"`
	Lon           float64       `json:"lon"`
	Hours         int8          `json:"hours"`
	Location      Location      `json:"location"`
	Weather       OpenMeteoResp `json:"weather"`
}

const cacheFile = "/tmp/weather_cache.json"
const cacheSchemaVersion = 1
const cacheTTL = 30 * time.Minute

func LoadCache(lat, lon float64, hours int8) (*CacheData, bool, error) {
	file, err := os.Open(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	defer file.Close()

	var data CacheData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, false, err
	}

	if data.SchemaVersion != cacheSchemaVersion {
		return nil, false, nil
	}
	if data.ExpiresAt.IsZero() || time.Now().After(data.ExpiresAt) {
		return nil, false, nil
	}
	if data.Hours != hours || !approxEqual(data.Lat, lat) || !approxEqual(data.Lon, lon) {
		return nil, false, nil
	}

	return &data, true, nil
}

func SaveCache(loc Location, weather OpenMeteoResp, lat, lon float64, hours int8) error {
	cacheDir := filepath.Dir(cacheFile)
	file, err := os.CreateTemp(cacheDir, "weather_cache_*.json")
	if err != nil {
		return err
	}

	tempName := file.Name()
	defer os.Remove(tempName)

	now := time.Now()
	data := CacheData{
		SchemaVersion: cacheSchemaVersion,
		CachedAt:      now,
		ExpiresAt:     now.Add(cacheTTL),
		Lat:           lat,
		Lon:           lon,
		Hours:         hours,
		Location:      loc,
		Weather:       weather,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	if err := os.Rename(tempName, cacheFile); err != nil {
		return err
	}

	return nil
}

func approxEqual(a, b float64) bool {
	return math.Abs(a-b) <= 0.0001
}
