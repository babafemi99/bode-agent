package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/babafemi99/bode-agent/internal/bot"
	"github.com/joho/godotenv"
	"github.com/logrusorgru/aurora/v3"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	fmt.Println(
		aurora.Bold(
			aurora.Magenta(banner),
		),
	)

	fmt.Printf(
		"  %-20s %s\n",
		aurora.Cyan("Service:"),
		aurora.White(ServiceName),
	)

	fmt.Printf(
		"  %-20s %s\n",
		aurora.Cyan("Transport:"),
		aurora.White("Telegram"),
	)

	fmt.Printf(
		"  %-20s %s\n",
		aurora.Cyan("Mode:"),
		aurora.White("Long Polling"),
	)

	fmt.Println()

	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env file not loaded: %v", err)
	}

	token := os.Getenv("KEY")
	if token == "" {
		fmt.Println(
			aurora.Red(
				aurora.Bold("  ✗ TELEGRAM_BOT_TOKEN is not set"),
			),
		)
		os.Exit(1)
	}

	fmt.Println(
		aurora.Yellow("  ◌ Starting Bọ̀dé Bot..."),
	)

	if err := bot.New(ctx, token); err != nil {
		fmt.Println(
			aurora.Red(
				aurora.Bold(
					fmt.Sprintf("  ✗ Failed to start bot: %v", err),
				),
			),
		)
		os.Exit(1)
	}

	fmt.Println(
		aurora.Green(
			aurora.Bold("  ✓ Telegram Bot       Connected"),
		),
	)

	fmt.Println(
		aurora.Green(
			aurora.Bold("  ✓ Update Polling     Active"),
		),
	)

	fmt.Println(
		aurora.Green(
			aurora.Bold("  ✓ Message Handler    Running"),
		),
	)

	fmt.Println()

	fmt.Println(
		aurora.Bold(
			aurora.Green("  ✓ Bọ̀dé Bot is ready. Let's go. 🚀"),
		),
	)

	fmt.Println()

	<-ctx.Done()

	fmt.Println()
	fmt.Println(
		aurora.Bold(
			aurora.Yellow("  ⚠ Shutdown signal received"),
		),
	)

	fmt.Println(
		aurora.Faint("  Stopping Bọ̀dé Bot..."),
	)

	time.Sleep(500 * time.Millisecond)

	fmt.Println(
		aurora.Green("  ✓ Bọ̀dé Bot stopped."),
	)
}
