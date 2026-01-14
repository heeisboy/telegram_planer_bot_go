package timeparse

import "time"

const Layout = "2006-01-02 15:04"

// ParseInLocation parses a local datetime string in the provided location.
func ParseInLocation(value string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.UTC
	}
	return time.ParseInLocation(Layout, value, loc)
}

// IsSameDay reports whether t is on the same calendar day as ref in ref's location.
func IsSameDay(t time.Time, ref time.Time, loc *time.Location) bool {
	if loc == nil {
		loc = time.UTC
	}
	t = t.In(loc)
	ref = ref.In(loc)
	return t.Year() == ref.Year() && t.Month() == ref.Month() && t.Day() == ref.Day()
}
