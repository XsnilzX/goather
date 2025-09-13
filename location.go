package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Ergebnisstruktur für den Aufrufer
type Location struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	City    string  `json:"city"`
	Region  string  `json:"region"`
	Country string  `json:"country"`
	Source  string  `json:"source"`
}

// Provider-Interface
type provider interface {
	Name() string
	Lookup(ctx context.Context, client *http.Client) (Location, error)
	qualityHint(loc Location) int // höhere Zahl = besser
}

// ====== Provider 1: ipapi.co ======
type ipapiCo struct{}

func (p ipapiCo) Name() string { return "ipapi.co" }

func (p ipapiCo) Lookup(ctx context.Context, client *http.Client) (Location, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://ipapi.co/json/", nil)
	req.Header.Set("User-Agent", "weather-widget/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Location{}, fmt.Errorf("ipapi.co status %d", resp.StatusCode)
	}
	var r struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		City      string  `json:"city"`
		Region    string  `json:"region"`
		Country   string  `json:"country_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Location{}, err
	}
	return Location{
		Lat:     r.Latitude,
		Lon:     r.Longitude,
		City:    r.City,
		Region:  r.Region,
		Country: r.Country,
		Source:  p.Name(),
	}, nil
}

func (p ipapiCo) qualityHint(loc Location) int { // city+region bekannt ist gut
	score := 0
	if loc.City != "" {
		score += 2
	}
	if loc.Region != "" {
		score += 2
	}
	if loc.Country != "" {
		score++
	}
	if loc.Lat != 0 || loc.Lon != 0 {
		score += 2
	}
	return score
}

// ====== Provider 2: ip-api.com ======
type ipapiCom struct{}

func (p ipapiCom) Name() string { return "ip-api.com" }

func (p ipapiCom) Lookup(ctx context.Context, client *http.Client) (Location, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://ip-api.com/json/?fields=status,country,regionName,city,lat,lon,message", nil)
	req.Header.Set("User-Agent", "weather-widget/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()
	var r struct {
		Status    string  `json:"status"`
		Message   string  `json:"message"`
		Country   string  `json:"country"`
		Region    string  `json:"regionName"`
		City      string  `json:"city"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Location{}, err
	}
	if r.Status != "success" {
		if r.Message == "" {
			r.Message = "unknown error"
		}
		return Location{}, errors.New(r.Message)
	}
	return Location{
		Lat:     r.Latitude,
		Lon:     r.Longitude,
		City:    r.City,
		Region:  r.Region,
		Country: r.Country,
		Source:  p.Name(),
	}, nil
}

func (p ipapiCom) qualityHint(loc Location) int {
	score := 0
	if loc.City != "" {
		score += 2
	}
	if loc.Region != "" {
		score += 2
	}
	if loc.Country != "" {
		score++
	}
	if loc.Lat != 0 || loc.Lon != 0 {
		score += 2
	}
	return score
}

// ====== Provider 3: ipwho.is ======
type ipwhoIs struct{}

func (p ipwhoIs) Name() string { return "ipwho.is" }

func (p ipwhoIs) Lookup(ctx context.Context, client *http.Client) (Location, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://ipwho.is/", nil)
	req.Header.Set("User-Agent", "weather-widget/1.0")
	resp, err := client.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()
	var r struct {
		Success   bool    `json:"success"`
		Message   string  `json:"message"`
		City      string  `json:"city"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Location{}, err
	}
	if !r.Success {
		if r.Message == "" {
			r.Message = "lookup failed"
		}
		return Location{}, errors.New(r.Message)
	}
	return Location{
		Lat:     r.Latitude,
		Lon:     r.Longitude,
		City:    r.City,
		Region:  r.Region,
		Country: r.Country,
		Source:  p.Name(),
	}, nil
}

func (p ipwhoIs) qualityHint(loc Location) int {
	score := 0
	if loc.City != "" {
		score += 2
	}
	if loc.Region != "" {
		score += 2
	}
	if loc.Country != "" {
		score++
	}
	if loc.Lat != 0 || loc.Lon != 0 {
		score += 2
	}
	return score
}

// ====== Public API: GetLocation ======

type Option func(*options)

type options struct {
	OverallTimeout time.Duration
	PerReqTimeout  time.Duration
	PreferFastest  bool // true: nimm erste brauchbare Antwort; false: sammle kurz & wähle beste
	ExtraProviders []provider
}

func WithOverallTimeout(d time.Duration) Option { return func(o *options) { o.OverallTimeout = d } }
func WithPerReqTimeout(d time.Duration) Option  { return func(o *options) { o.PerReqTimeout = d } }
func WithPreferFastest(v bool) Option           { return func(o *options) { o.PreferFastest = v } }
func WithExtraProviders(ps ...provider) Option {
	return func(o *options) { o.ExtraProviders = append(o.ExtraProviders, ps...) }
}

// GetLocation fragt mehrere Provider parallel ab.
func GetLocation(ctx context.Context, opts ...Option) (Location, error) {
	// Defaults
	cfg := &options{
		OverallTimeout: 1500 * time.Millisecond,
		PerReqTimeout:  900 * time.Millisecond,
		PreferFastest:  true,
	}
	for _, o := range opts {
		o(cfg)
	}

	providers := []provider{
		ipapiCo{},
		ipapiCom{},
		ipwhoIs{},
	}
	providers = append(providers, cfg.ExtraProviders...)

	type result struct {
		loc Location
		err error
		dt  time.Duration
		p   provider
	}

	ctxOverall, cancel := context.WithTimeout(ctx, cfg.OverallTimeout)
	defer cancel()

	client := &http.Client{Timeout: cfg.PerReqTimeout}
	results := make(chan result, len(providers))

	// Fan-out
	for _, p := range providers {
		go func(p provider) {
			start := time.Now()
			ctxReq, cancelReq := context.WithTimeout(ctxOverall, cfg.PerReqTimeout)
			defer cancelReq()
			loc, err := p.Lookup(ctxReq, client)
			results <- result{loc: loc, err: err, dt: time.Since(start), p: p}
		}(p)
	}

	var (
		best     *result
		gotAny   bool
		fastestT time.Duration
	)

	// Wenn PreferFastest=true: nimm erste brauchbare sofort.
	// Sonst: sammle bis OverallTimeout und nimm best-qualitativ.
	for i := 0; i < len(providers); i++ {
		select {
		case r := <-results:
			if r.err == nil {
				gotAny = true
				// erste brauchbare?
				if cfg.PreferFastest && best == nil {
					return r.loc, nil
				}
				// wir bewerten Qualität
				q := r.p.qualityHint(r.loc)
				if best == nil || q > r.p.qualityHint(best.loc) || (q == r.p.qualityHint(best.loc) && r.dt < fastestT) {
					tmp := r // copy
					best = &tmp
					fastestT = r.dt
				}
			}
		case <-ctxOverall.Done():
			// Timeout – beende früh
			if cfg.PreferFastest {
				// evtl. noch kein Resultat -> break outer
				i = len(providers) // break for
			} else {
				// wir geben zurück, was wir haben
				i = len(providers)
			}
		}
	}

	if best != nil {
		return best.loc, nil
	}
	if gotAny {
		// theoretisch nie, aber zur Sicherheit
		return Location{}, errors.New("no usable location result selected")
	}
	return Location{}, errors.New("all location providers failed or timed out")
}
