# goather

[![Flake Check](https://github.com/XsnilzX/goather/actions/workflows/flake-check.yml/badge.svg)](https://github.com/XsnilzX/goather/actions/workflows/flake-check.yml)

A small Go-based weather widget that prints JSON for Waybar. It uses IP-based
geolocation, fetches weather from Open-Meteo, and keeps a short cache in
`${XDG_CACHE_HOME:-~/.cache}/goather/weather.json`.

## Features

- JSON output with `text`, `tooltip`, and `class` fields
- IP geolocation with concurrent providers
- Hourly forecast snippet in the tooltip
- Simple cache to reduce API calls

## Install with Nix

### Run directly

```sh
nix run .
```

### Build the package

```sh
nix build .
./result/bin/goather
```

### Install in a profile

```sh
nix profile install .
goather
```

### NixOS module

```nix
{
  inputs.goather.url = "github:XsnilzX/goather";

  outputs = { self, nixpkgs, goather, ... }: {
    nixosConfigurations.myHost = nixpkgs.lib.nixosSystem {
      modules = [
        goather.nixosModules.default
        {
          programs.goather.enable = true;
        }
      ];
    };
  };
}
```

### Home Manager module

```nix
{
  inputs.goather.url = "github:XsnilzX/goather";

  outputs = { self, nixpkgs, goather, ... }: {
    homeConfigurations.me = nixpkgs.lib.homeManagerConfiguration {
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
      modules = [
        goather.homeManagerModules.default
        {
          programs.goather.enable = true;
        }
      ];
    };
  };
}
```

## Usage

Run the binary and capture its JSON output:

```sh
goather
```

Example output:

```json
{
  "text": "☀️ 21°C",
  "tooltip": "City, Country\nKlarer Himmel\n🌡️ Temperature: 21°C (feels 20°C)\n💧 Humidity: 40%\n💨 Wind: 12 km/h\n⏰ Next 6 hours:\n14:00 ☀️ 21°C\n15:00 🌤️ 22°C\nUpdated: 14:05",
  "class": "clear"
}
```

### Waybar example

```json
"custom/weather": {
  "exec": "goather",
  "return-type": "json",
  "interval": 1800,
  "format": "{text}",
  "tooltip": true
}
```

`goather` is a one-shot CLI. Waybar should handle refreshes via `interval` or a
signal-driven setup.

Geolocation is IP-based, so VPNs, proxies, or privacy relays can affect the
detected location.

## Build from source (non-Nix)

```sh
make
./goather
```

## Development

- Go version: 1.26.1
- Format: `gofmt -w .`
- Tests: `go test ./...`

## CI

- Flake checks run on pull requests and pushes to `main`/`master`.
- Full multi-system validation: `nix flake check --all-systems`.

## License

MIT
