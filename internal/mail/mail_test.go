package mail

import (
	"net/smtp"
	"strings"
	"testing"

	"github.com/kristyancarvalho/tux-letter/internal/config"
)

func TestBuildMultipart(t *testing.T) {
	raw := string(Build(Message{
		From:    "from@example.com",
		To:      []string{"a@example.com", "b@example.com"},
		Subject: "Tux Letter",
		Text:    "plain body",
		HTML:    "<p>html body</p>",
	}))

	for _, want := range []string{
		"From: from@example.com",
		"To: a@example.com, b@example.com",
		"Subject: Tux Letter",
		"multipart/alternative",
		"text/plain; charset=UTF-8",
		"text/html; charset=UTF-8",
		"plain body",
		"<p>html body</p>",
	} {
		if !strings.Contains(raw, want) {
			t.Errorf("built message missing %q", want)
		}
	}
}

func TestBuildSanitizesSubjectHeader(t *testing.T) {
	raw := string(Build(Message{Subject: "Evil\r\nBcc: x@y.com", Text: "t"}))
	if strings.Contains(raw, "\r\nBcc:") {
		t.Errorf("header injection not neutralized: %q", raw)
	}
	if !strings.Contains(raw, "Subject: Evil Bcc: x@y.com") {
		t.Errorf("expected folded subject, got %q", raw)
	}
}

func TestResolveReadsEnv(t *testing.T) {
	cfg := config.Default().Email
	t.Setenv(cfg.SMTPHostEnv, "smtp.example.com")
	t.Setenv(cfg.FromEnv, "from@example.com")
	t.Setenv(cfg.ToEnv, "a@example.com, b@example.com")

	s, err := Resolve(cfg)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if s.Host != "smtp.example.com" {
		t.Errorf("host: %q", s.Host)
	}
	if s.Port != "587" {
		t.Errorf("expected default port 587, got %q", s.Port)
	}
	if len(s.To) != 2 {
		t.Errorf("expected 2 recipients, got %v", s.To)
	}
}

func TestResolveMissingRequired(t *testing.T) {
	cfg := config.Default().Email
	t.Setenv(cfg.SMTPHostEnv, "")
	t.Setenv(cfg.FromEnv, "")
	t.Setenv(cfg.ToEnv, "")
	if _, err := Resolve(cfg); err == nil {
		t.Fatal("expected error for missing env vars")
	}
}

func TestSendUsesResolvedRecipients(t *testing.T) {
	var gotAddr, gotFrom string
	var gotTo []string
	sender := func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		gotAddr, gotFrom, gotTo = addr, from, to
		return nil
	}
	s := Settings{Host: "smtp.example.com", Port: "587", From: "f@example.com", To: []string{"t@example.com"}}
	if err := send(sender, s, Message{Subject: "Hi", Text: "x"}); err != nil {
		t.Fatalf("send: %v", err)
	}
	if gotAddr != "smtp.example.com:587" {
		t.Errorf("addr: %q", gotAddr)
	}
	if gotFrom != "f@example.com" {
		t.Errorf("from: %q", gotFrom)
	}
	if len(gotTo) != 1 || gotTo[0] != "t@example.com" {
		t.Errorf("to: %v", gotTo)
	}
}
