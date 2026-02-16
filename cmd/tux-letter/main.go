package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kristyancarvalho/tux-letter/internal/ai"
	"github.com/kristyancarvalho/tux-letter/internal/database"
	"github.com/kristyancarvalho/tux-letter/internal/email"
	"github.com/kristyancarvalho/tux-letter/internal/scraper"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load()

	dbPath := getEnv("DB_PATH", "./data/tux-letter.db")
	if err := database.Init(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	cronSchedule := getEnv("CRON_SCHEDULE", "0 20 * * *")

	c := cron.New()
	c.AddFunc(cronSchedule, runJob)

	fmt.Printf("Tux Letter started. Schedule: %s\n", cronSchedule)
	fmt.Println("Press Ctrl+C to stop")

	runJob()

	c.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down gracefully...")
	c.Stop()
}

func runJob() {
	fmt.Println("\n=== Starting scraping job ===")

	news, err := scraper.ScrapeAll("sites.json")
	if err != nil {
		log.Printf("Scraping error: %v", err)
		return
	}

	if len(news) == 0 {
		fmt.Println("No new articles found")
		return
	}

	fmt.Printf("Found %d new articles\n", len(news))

	fmt.Println("Synthesizing news with AI...")
	synthesized, err := ai.SynthesizeNews(news)
	if err != nil {
		log.Printf("AI synthesis error: %v", err)
		return
	}

	fmt.Println("Sending email...")
	if err := email.SendNews(synthesized); err != nil {
		log.Printf("Email error: %v", err)
		return
	}

	fmt.Println("✅ Email sent successfully")
	fmt.Println("=== Job completed ===")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
