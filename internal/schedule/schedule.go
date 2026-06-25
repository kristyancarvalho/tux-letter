package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseTimeOfDay(s string) (int, int, error) {
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid schedule %q, want HH:MM", s)
	}
	hour, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("invalid hour in schedule %q", s)
	}
	min, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || min < 0 || min > 59 {
		return 0, 0, fmt.Errorf("invalid minute in schedule %q", s)
	}
	return hour, min, nil
}

func NextRun(now time.Time, hhmm string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc = time.UTC
	}
	hour, min, err := ParseTimeOfDay(hhmm)
	if err != nil {
		return time.Time{}, err
	}
	local := now.In(loc)
	next := time.Date(local.Year(), local.Month(), local.Day(), hour, min, 0, 0, loc)
	if !next.After(local) {
		next = next.AddDate(0, 0, 1)
	}
	return next, nil
}

func LoadLocation(name string) (*time.Location, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return time.UTC, nil
	}
	return time.LoadLocation(name)
}
