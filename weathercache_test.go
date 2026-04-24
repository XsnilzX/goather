package main

import (
	"path/filepath"
	"testing"
)

func TestCacheFilePathUsesXDGCacheHome(t *testing.T) {
	cacheHome := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cacheHome)

	got := cacheFilePath()
	want := filepath.Join(cacheHome, "goather", "weather.json")
	if got != want {
		t.Fatalf("expected cache path %q, got %q", want, got)
	}
}

func TestSaveAndLoadCacheRoundTrip(t *testing.T) {
	cacheHome := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cacheHome)

	loc := Location{
		Lat:     52.52,
		Lon:     13.41,
		City:    "Berlin",
		Country: "Germany",
	}
	weather := OpenMeteoResp{}
	weather.Current.Temperature2m = 21
	weather.Current.WeatherCode = 0

	if err := SaveCache(loc, weather, loc.Lat, loc.Lon, 6); err != nil {
		t.Fatalf("save cache: %v", err)
	}

	cached, hit, err := LoadCache(loc.Lat, loc.Lon, 6)
	if err != nil {
		t.Fatalf("load cache: %v", err)
	}
	if !hit {
		t.Fatal("expected cache hit after saving cache")
	}
	if cached == nil {
		t.Fatal("expected cached data")
	}
	if cached.Location.City != loc.City {
		t.Fatalf("expected cached city %q, got %q", loc.City, cached.Location.City)
	}
}
