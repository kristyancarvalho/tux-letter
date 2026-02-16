package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kristyancarvalho/tux-letter/internal/scraper"
)

type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type SynthesizedNews struct {
	Summary    string
	References []scraper.NewsItem
}

func SynthesizeNews(news []scraper.NewsItem) (*SynthesizedNews, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	model := os.Getenv("OPENROUTER_MODEL")

	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY not set")
	}

	if model == "" {
		model = "anthropic/claude-3.5-sonnet"
	}

	prompt := buildPrompt(news)

	reqBody := OpenRouterRequest{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("OpenRouter API error: %s - %s", resp.Status, string(body))
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return nil, err
	}

	if len(openRouterResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenRouter")
	}

	summary := openRouterResp.Choices[0].Message.Content

	return &SynthesizedNews{
		Summary:    summary,
		References: news,
	}, nil
}

func buildPrompt(news []scraper.NewsItem) string {
	var builder strings.Builder

	builder.WriteString("Você é um especialista em Linux e Open Source. Crie uma newsletter profissional em português brasileiro.\n\n")
	builder.WriteString("FORMATO OBRIGATÓRIO:\n\n")
	builder.WriteString("1. PRIMEIRO PARÁGRAFO: Visão geral de 2-3 frases sobre os principais destaques do dia\n\n")
	builder.WriteString("2. TÓPICOS TEMÁTICOS: Organize as notícias por temas usando #### para títulos\n")
	builder.WriteString("   - Use títulos claros: #### Kernel e Sistema, #### Distribuições, #### Segurança, etc\n")
	builder.WriteString("   - Escreva em parágrafos corridos (NÃO use listas ou bullet points)\n")
	builder.WriteString("   - Adicione [N] após mencionar cada notícia (onde N é o número da referência)\n")
	builder.WriteString("   - Use **negrito** apenas para nomes de projetos/ferramentas importantes\n\n")
	builder.WriteString("NOTÍCIAS:\n\n")

	for i, item := range news {
		builder.WriteString(fmt.Sprintf("[%d] %s - %s\n", i+1, item.Site, item.Title))
	}

	builder.WriteString("\n\nEXEMPLO DE FORMATO:\n\n")
	builder.WriteString("Esta edição traz atualizações importantes no kernel Linux, novidades em distribuições populares e debates sobre privacidade na comunidade open source.\n\n")
	builder.WriteString("#### Kernel e Sistema\n\n")
	builder.WriteString("O kernel Linux recebeu patches críticos de segurança [1] que corrigem vulnerabilidades descobertas recentemente. Paralelamente, a comunidade debate mudanças propostas para o scheduler [3], buscando melhor performance em sistemas multi-core modernos.\n\n")
	builder.WriteString("#### Distribuições\n\n")
	builder.WriteString("O **Ubuntu 24.04 LTS** anunciou sua data de lançamento [2], trazendo melhorias significativas no desempenho e suporte estendido. Enquanto isso, o **Fedora** trabalha em integrações mais profundas com Wayland [5].\n\n")

	return builder.String()
}
