package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/kristyancarvalho/tux-letter/internal/ai"
)

type EmailConfig struct {
	SMTPHost string
	SMTPPort string
	User     string
	Pass     string
	From     string
	To       string
}

func LoadConfig() EmailConfig {
	return EmailConfig{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),
		User:     os.Getenv("SMTP_USER"),
		Pass:     os.Getenv("SMTP_PASS"),
		From:     os.Getenv("EMAIL_FROM"),
		To:       os.Getenv("EMAIL_TO"),
	}
}

func SendNews(synthesized *ai.SynthesizedNews) error {
	config := LoadConfig()

	htmlBody := generateHTML(synthesized)

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: Tux Letter - Daily Linux News\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		config.From,
		config.To,
		htmlBody,
	))

	auth := smtp.PlainAuth("", config.User, config.Pass, config.SMTPHost)
	addr := config.SMTPHost + ":" + config.SMTPPort

	return smtp.SendMail(addr, auth, config.From, []string{config.To}, msg)
}

func generateHTML(synthesized *ai.SynthesizedNews) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { 
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
			background: #0d1117;
			color: #e6edf3;
			padding: 20px;
			line-height: 1.7;
		}
		.container { 
			max-width: 700px;
			margin: 0 auto;
			background: #161b22;
			border-radius: 8px;
			overflow: hidden;
			border: 1px solid #30363d;
		}
		.header {
			background: linear-gradient(135deg, #238636 0%, #2ea043 100%);
			padding: 30px;
			text-align: center;
		}
		.header h1 {
			color: white;
			font-size: 1.8em;
			font-weight: 600;
			margin-bottom: 5px;
		}
		.header .date {
			color: rgba(255,255,255,0.85);
			font-size: 0.85em;
		}
		.divider {
			height: 1px;
			background: #30363d;
			margin: 0;
		}
		.section {
			padding: 30px;
		}
		.section-title {
			color: #7ee787;
			font-size: 1.1em;
			font-weight: 600;
			margin-bottom: 15px;
			padding-bottom: 8px;
			border-bottom: 2px solid #238636;
		}
		.overview {
			background: #0d1117;
			padding: 20px;
			border-radius: 6px;
			border-left: 3px solid #58a6ff;
			color: #c9d1d9;
			font-size: 0.95em;
			line-height: 1.6;
		}
		.content h3 {
			color: #58a6ff;
			font-size: 1.05em;
			font-weight: 600;
			margin: 25px 0 12px 0;
			padding-left: 12px;
			border-left: 3px solid #58a6ff;
		}
		.content p {
			color: #c9d1d9;
			margin-bottom: 14px;
			text-align: justify;
		}
		.content strong {
			color: #e6edf3;
			font-weight: 600;
		}
		.ref-link {
			display: inline-block;
			background: #1f6feb;
			color: white;
			padding: 1px 6px;
			border-radius: 3px;
			font-size: 0.8em;
			font-weight: 600;
			text-decoration: none;
			margin: 0 2px;
		}
		.references {
			background: #0d1117;
			padding: 25px;
		}
		.ref-title {
			color: #7ee787;
			font-size: 1.1em;
			font-weight: 600;
			margin-bottom: 15px;
		}
		.ref-item {
			background: #161b22;
			padding: 12px;
			margin-bottom: 8px;
			border-radius: 4px;
			border-left: 3px solid #238636;
			font-size: 0.9em;
		}
		.ref-item:hover {
			background: #1c2128;
		}
		.ref-num {
			display: inline-block;
			background: #1f6feb;
			color: white;
			min-width: 24px;
			height: 24px;
			border-radius: 12px;
			text-align: center;
			line-height: 24px;
			font-weight: 600;
			font-size: 0.85em;
			margin-right: 10px;
		}
		.ref-site {
			display: inline-block;
			background: #238636;
			color: white;
			padding: 2px 8px;
			border-radius: 10px;
			font-size: 0.75em;
			font-weight: 600;
			margin-left: 8px;
		}
		.ref-item a {
			color: #58a6ff;
			text-decoration: none;
		}
		.ref-item a:hover {
			text-decoration: underline;
		}
		.footer {
			padding: 20px;
			text-align: center;
			color: #8b949e;
			font-size: 0.8em;
			border-top: 1px solid #30363d;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>🐧 Tux Letter</h1>
			<div class="date">{{.Date}}</div>
		</div>

		<div class="divider"></div>

		<div class="section">
			<div class="section-title">📋 Visão Geral</div>
			<div class="overview">
				{{.Overview}}
			</div>
		</div>

		<div class="divider"></div>

		<div class="section">
			<div class="section-title">📰 Análise Detalhada</div>
			<div class="content">
				{{.Content}}
			</div>
		</div>

		<div class="divider"></div>

		<div class="references">
			<div class="ref-title">📚 Referências ({{.Total}} notícias)</div>
			{{range .References}}
			<div class="ref-item">
				<span class="ref-num">{{.Index}}</span>
				<a href="{{.Link}}" target="_blank">{{.Title}}</a>
				<span class="ref-site">{{.Site}}</span>
			</div>
			{{end}}
		</div>

		<div class="footer">
			Tux Letter • Newsletter automática sobre Linux e Open Source.
		</div>
	</div>
</body>
</html>
`

	t := template.Must(template.New("email").Parse(tmpl))

	contentHTML, overview := formatContent(synthesized.Summary, len(synthesized.References))

	data := struct {
		Date       string
		Overview   template.HTML
		Content    template.HTML
		Total      int
		References []struct {
			Index int
			Site  string
			Title string
			Link  string
		}
	}{
		Date:     formatDate(),
		Overview: template.HTML(overview),
		Content:  template.HTML(contentHTML),
		Total:    len(synthesized.References),
	}

	for i, ref := range synthesized.References {
		data.References = append(data.References, struct {
			Index int
			Site  string
			Title string
			Link  string
		}{
			Index: i + 1,
			Site:  ref.Site,
			Title: ref.Title,
			Link:  ref.Link,
		})
	}

	var buf bytes.Buffer
	t.Execute(&buf, data)

	return buf.String()
}

func formatDate() string {
	now := time.Now()
	months := []string{
		"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
		"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
	}
	weekdays := []string{
		"Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado",
	}
	return fmt.Sprintf("%s, %d de %s de %d",
		weekdays[now.Weekday()],
		now.Day(),
		months[now.Month()-1],
		now.Year())
}

func formatContent(markdown string, totalNews int) (string, string) {
	lines := strings.Split(markdown, "\n")

	var overviewLines []string
	var contentLines []string
	inOverview := true

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "####") {
			inOverview = false
		}

		if inOverview && trimmed != "" && !strings.HasPrefix(trimmed, "#") {
			overviewLines = append(overviewLines, trimmed)
		} else if !inOverview {
			contentLines = append(contentLines, line)
		}
	}

	overview := strings.Join(overviewLines, " ")
	if overview == "" {
		overview = fmt.Sprintf("Nesta edição: %d novas notícias sobre Linux e Open Source, incluindo atualizações de kernel, distribuições, ferramentas e tendências da comunidade.", totalNews)
	}

	content := strings.Join(contentLines, "\n")

	content = strings.ReplaceAll(content, "####", "<h3>")
	content = strings.ReplaceAll(content, "**", "<strong>")

	for i := 1; i <= 100; i++ {
		ref := fmt.Sprintf("[%d]", i)
		link := fmt.Sprintf("<a href=\"#ref-%d\" class=\"ref-link\">%d</a>", i, i)
		content = strings.ReplaceAll(content, ref, link)
	}

	paragraphs := strings.Split(content, "\n\n")
	var formatted []string
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p != "" && !strings.HasPrefix(p, "<h3>") {
			formatted = append(formatted, "<p>"+p+"</p>")
		} else if p != "" {
			formatted = append(formatted, p)
		}
	}
	content = strings.Join(formatted, "\n")

	return content, overview
}
