package schedule

import (
	"testing"
	"time"
)

func TestParseTimeOfDay(t *testing.T) {
	h, m, err := ParseTimeOfDay("20:05")
	if err != nil || h != 20 || m != 5 {
		t.Fatalf("got %d:%d err=%v", h, m, err)
	}
	for _, bad := range []string{"", "20", "24:00", "20:60", "aa:bb", "-1:00"} {
		if _, _, err := ParseTimeOfDay(bad); err == nil {
			t.Errorf("expected error for %q", bad)
		}
	}
}

func TestNextRunLaterToday(t *testing.T) {
	loc := time.UTC
	now := time.Date(2026, 6, 25, 8, 0, 0, 0, loc)
	next, err := NextRun(now, "20:00", loc)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 6, 25, 20, 0, 0, 0, loc)
	if !next.Equal(want) {
		t.Errorf("got %v want %v", next, want)
	}
}

func TestNextRunRollsToTomorrow(t *testing.T) {
	loc := time.UTC
	now := time.Date(2026, 6, 25, 21, 0, 0, 0, loc)
	next, err := NextRun(now, "20:00", loc)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 6, 26, 20, 0, 0, 0, loc)
	if !next.Equal(want) {
		t.Errorf("got %v want %v", next, want)
	}
}

func TestNextRunRespectsTimezone(t *testing.T) {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		t.Skip("tzdata unavailable")
	}
	now := time.Date(2026, 6, 25, 12, 0, 0, 0, time.UTC)
	next, err := NextRun(now, "20:00", loc)
	if err != nil {
		t.Fatal(err)
	}
	h, m, _ := next.Clock()
	if h != 20 || m != 0 {
		t.Errorf("expected 20:00 local, got %02d:%02d", h, m)
	}
	if next.Location().String() != "America/Sao_Paulo" {
		t.Errorf("expected Sao_Paulo location, got %s", next.Location())
	}
}

func TestNextRunInvalidSchedule(t *testing.T) {
	if _, err := NextRun(time.Now(), "nope", time.UTC); err == nil {
		t.Error("expected error for invalid schedule")
	}
}
