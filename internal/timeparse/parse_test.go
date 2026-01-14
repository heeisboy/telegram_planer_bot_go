package timeparse

import (
	"testing"
	"time"
)

func TestParseInLocation(t *testing.T) {
	loc := time.FixedZone("Test", 2*60*60)
	parsed, err := ParseInLocation("2024-10-01 09:30", loc)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.Location() != loc {
		t.Fatalf("expected location to be preserved")
	}
	if parsed.Year() != 2024 || parsed.Month() != 10 || parsed.Day() != 1 || parsed.Hour() != 9 || parsed.Minute() != 30 {
		t.Fatalf("unexpected datetime: %v", parsed)
	}
}

func TestIsSameDay(t *testing.T) {
	loc := time.FixedZone("Test", -5*60*60)
	t1 := time.Date(2024, 10, 1, 23, 0, 0, 0, loc)
	t2 := time.Date(2024, 10, 1, 1, 0, 0, 0, loc)
	if !IsSameDay(t1, t2, loc) {
		t.Fatalf("expected same day")
	}
	t3 := time.Date(2024, 10, 2, 0, 30, 0, 0, loc)
	if IsSameDay(t1, t3, loc) {
		t.Fatalf("expected different day")
	}
}
