APP_NAME = goather
CACHE_FILE = /tmp/weather_cache.json

.PHONY: all debug release clean

# Standard: Release-Build
all: release

# Debug-Build mit Symbolen
debug:
	go build -o $(APP_NAME) .

# Optimierter Release-Build (klein & schnell)
release:
	go build -trimpath -ldflags="-s -w" -o $(APP_NAME) .

# Aufräumen (Binary + Cache-Datei)
clean:
	rm -f $(APP_NAME)
	rm -f $(CACHE_FILE)

