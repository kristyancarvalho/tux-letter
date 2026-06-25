package mail

import (
	"fmt"
	"strings"
)

const boundary = "tux-letter-boundary-2f8a1c9e4b7d"

type Message struct {
	From    string
	To      []string
	Subject string
	Text    string
	HTML    string
}

func Build(m Message) []byte {
	var b strings.Builder
	b.WriteString("From: " + m.From + "\r\n")
	b.WriteString("To: " + strings.Join(m.To, ", ") + "\r\n")
	b.WriteString("Subject: " + sanitizeHeader(m.Subject) + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q\r\n", boundary))
	b.WriteString("\r\n")

	if m.Text != "" {
		b.WriteString("--" + boundary + "\r\n")
		b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		b.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		b.WriteString(normalizeBody(m.Text))
		b.WriteString("\r\n")
	}

	if m.HTML != "" {
		b.WriteString("--" + boundary + "\r\n")
		b.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		b.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
		b.WriteString(normalizeBody(m.HTML))
		b.WriteString("\r\n")
	}

	b.WriteString("--" + boundary + "--\r\n")
	return []byte(b.String())
}

func sanitizeHeader(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func normalizeBody(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\n", "\r\n")
}
