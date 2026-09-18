package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/babafemi99/bode-agent/internal/bot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env file not loaded: %v", err)
	}

	token := os.Getenv("KEY")
	if token == "" {
		log.Fatal("KEY is not set")
	}

	mode := os.Getenv("BOT_MODE")
	if mode == "" {
		log.Fatal("BOT_MODE is not set")
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
	if err != nil {
		log.Fatalf("invalid PORT: %v", err)
	}

	webhookURL := os.Getenv("WEBHOOK_URL")

	if mode == "webhook" && webhookURL == "" {
		log.Fatal("WEBHOOK_URL is not set")
	}

	b, err := bot.New(ctx, token, mode)
	if err != nil {
		log.Fatalf("failed to start bot: %v", err)
	}

	api := &API{
		Port:           port,
		Bot:            b,
		WebhookHandler: bot.NewWebhookHandler(b),
	}

	if mode == "webhook" {
		wh, err := tgbotapi.NewWebhook(
			webhookURL + "/webhook",
		)
		if err != nil {
			log.Fatalf("failed to create webhook: %v", err)
		}

		_, err = b.Client.Request(wh)
		if err != nil {
			log.Fatalf("failed to set webhook: %v", err)
		}

		log.Printf(
			"telegram webhook registered: %s/webhook",
			webhookURL,
		)
	}

	go func() {
		if err := api.Serve(); err != nil {
			log.Printf("API server stopped: %v", err)
			stop()
		}
	}()

	fmt.Println("Bọ̀Dé Bot is running.")

	<-ctx.Done()

	fmt.Println()
	fmt.Println("  [!] Shutdown signal received")
	fmt.Println("  Stopping Bọ̀Dé Bot...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := api.Shutdown(shutdownCtx); err != nil {
		log.Printf("API shutdown error: %v", err)
	}

	fmt.Println("  [+] Bọ̀Dé Bot stopped.")
}
