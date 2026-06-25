package urlx

import "testing"

func TestCanonical(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"https://example.com", "https://example.com/"},
		{"https://Example.com/Path/", "https://example.com/Path"},
		{"example.com", "https://example.com/"},
		{"https://example.com/a?utm_source=x&id=5", "https://example.com/a?id=5"},
		{"https://example.com/a#section", "https://example.com/a"},
		{"https://example.com:443/a", "https://example.com/a"},
		{"https://example.com/a?b=2&a=1", "https://example.com/a?a=1&b=2"},
	}
	for _, tt := range tests {
		got, err := Canonical(tt.in)
		if err != nil {
			t.Errorf("Canonical(%q) error: %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("Canonical(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestCanonicalErrors(t *testing.T) {
	for _, in := range []string{"", "   ", "ftp://example.com", "https://"} {
		if _, err := Canonical(in); err == nil {
			t.Errorf("Canonical(%q) expected error", in)
		}
	}
}

func TestHost(t *testing.T) {
	tests := map[string]string{
		"https://www.linux.com/news": "linux.com",
		"https://9to5linux.com":      "9to5linux.com",
		"archlinux.org/news":         "archlinux.org",
	}
	for in, want := range tests {
		got, err := Host(in)
		if err != nil {
			t.Errorf("Host(%q) error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("Host(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestResolve(t *testing.T) {
	got, err := Resolve("https://example.com/blog/", "../article")
	if err != nil {
		t.Fatal(err)
	}
	if got != "https://example.com/article" {
		t.Errorf("Resolve = %q", got)
	}
}
