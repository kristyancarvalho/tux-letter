package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/config"
)

type Settings struct {
	Host string
	Port string
	User string
	Pass string
	From string
	To   []string
}

func (s Settings) Addr() string {
	return s.Host + ":" + s.Port
}

func Resolve(cfg config.EmailConfig) (Settings, error) {
	s := Settings{
		Host: os.Getenv(cfg.SMTPHostEnv),
		Port: os.Getenv(cfg.SMTPPortEnv),
		User: os.Getenv(cfg.SMTPUserEnv),
		Pass: os.Getenv(cfg.SMTPPassEnv),
		From: os.Getenv(cfg.FromEnv),
		To:   splitRecipients(os.Getenv(cfg.ToEnv)),
	}
	if s.Port == "" {
		s.Port = "587"
	}

	var missing []string
	if s.Host == "" {
		missing = append(missing, cfg.SMTPHostEnv)
	}
	if s.From == "" {
		missing = append(missing, cfg.FromEnv)
	}
	if len(s.To) == 0 {
		missing = append(missing, cfg.ToEnv)
	}
	if len(missing) > 0 {
		return Settings{}, fmt.Errorf("missing required email environment variables: %s", strings.Join(missing, ", "))
	}
	return s, nil
}

func splitRecipients(raw string) []string {
	var out []string
	for _, part := range strings.FieldsFunc(raw, func(r rune) bool { return r == ',' || r == ';' }) {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

type SendFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

func Send(s Settings, m Message) error {
	return send(smtp.SendMail, s, m)
}

func send(sender SendFunc, s Settings, m Message) error {
	m.From = s.From
	m.To = s.To
	raw := Build(m)

	var auth smtp.Auth
	if s.User != "" && s.Pass != "" {
		auth = smtp.PlainAuth("", s.User, s.Pass, s.Host)
	}
	return sender(s.Addr(), auth, s.From, s.To, raw)
}
