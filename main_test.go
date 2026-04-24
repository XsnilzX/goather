package main

import (
	"encoding/json"
	"testing"
)

func TestWaybarOutputJSON(t *testing.T) {
	t.Parallel()

	out := WaybarOutput{
		Text:    "☀️ 21°C",
		Tooltip: "Weather tooltip",
		Class:   "clear",
	}

	data, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("marshal output: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	if payload["text"] != out.Text {
		t.Fatalf("expected text field %q, got %#v", out.Text, payload["text"])
	}
	if payload["tooltip"] != out.Tooltip {
		t.Fatalf("expected tooltip field %q, got %#v", out.Tooltip, payload["tooltip"])
	}
	if payload["class"] != out.Class {
		t.Fatalf("expected class field %q, got %#v", out.Class, payload["class"])
	}
	if _, ok := payload["display"]; ok {
		t.Fatal("unexpected legacy display field in waybar JSON")
	}
}
