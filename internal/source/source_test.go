package source

import "testing"

func TestNameFromHost(t *testing.T) {
	tests := map[string]string{
		"archlinux.org": "Archlinux",
		"www.linux.com": "Linux",
		"9to5linux.com": "9to5linux",
		"itsfoss.com":   "Itsfoss",
		"lwn.net":       "Lwn",
	}
	for in, want := range tests {
		if got := NameFromHost(in); got != want {
			t.Errorf("NameFromHost(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestBuildNormalizesAndDedupes(t *testing.T) {
	sources, errs := Build([]string{
		"https://9to5linux.com",
		"https://9to5linux.com/",
		"http://www.linux.com/news/",
	})
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(sources) != 2 {
		t.Fatalf("got %d sources, want 2 (dedupe failed): %+v", len(sources), sources)
	}
	if sources[0].Name != "9to5linux" {
		t.Errorf("name = %q, want 9to5linux", sources[0].Name)
	}
}

func TestBuildSkipsInvalidWithoutFailing(t *testing.T) {
	sources, errs := Build([]string{
		"https://valid.com",
		"ftp://invalid.com",
		"",
	})
	if len(sources) != 1 {
		t.Errorf("got %d valid sources, want 1", len(sources))
	}
	if len(errs) != 2 {
		t.Errorf("got %d errors, want 2", len(errs))
	}
}
