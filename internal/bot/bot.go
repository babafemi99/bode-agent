package bot

import (
	"context"
	"fmt"

	"github.com/babafemi99/bode-agent/pkg/lid"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	BotID int64
	ID    string

	Token  string
	Client *tgbotapi.BotAPI

	inbound  chan Message
	outbound chan Message
}

func New(ctx context.Context, token string) error {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return fmt.Errorf("failed to get Bot: %w", err)
	}

	client.Debug = true

	bot := &Bot{
		BotID:    client.Self.ID,
		ID:       lid.NewBot(),
		Token:    token,
		Client:   client,
		inbound:  make(chan Message, 256),
		outbound: make(chan Message, 256),
	}

	commands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Start using Bọ̀dé",
		},
		tgbotapi.BotCommand{
			Command:     "help",
			Description: "Show available commands",
		},
		tgbotapi.BotCommand{
			Command:     "events",
			Description: "View your events",
		},
		tgbotapi.BotCommand{
			Command:     "context",
			Description: "Show current context",
		},
	)

	_, err = bot.Client.Request(commands)
	if err != nil {
		return fmt.Errorf("failed to request commands: %w", err)
	}

	go bot.receive(ctx)
	go bot.send(ctx)
	go bot.handle(ctx)

	return nil
}

func (b *Bot) receive(ctx context.Context) {
	updates := b.Client.GetUpdatesChan(
		tgbotapi.UpdateConfig{
			Timeout: 30,
		},
	)

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				return
			}

			if update.Message == nil {
				continue
			}

			message := Message{
				BotID:     b.ID,
				UserID:    update.Message.From.ID,
				ChatID:    update.Message.Chat.ID,
				Text:      update.Message.Text,
				IsCommand: update.Message.IsCommand(),
				Command:   update.Message.Command(),
				Arguments: update.Message.CommandArguments(),
			}

			select {
			case b.inbound <- message:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (b *Bot) send(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case response := <-b.outbound:
			msg := tgbotapi.NewMessage(
				response.ChatID,
				response.Text,
			)

			_, err := b.Client.Send(msg)
			if err != nil {
				// handle error prolly just log for now
			}
		}
	}
}

func (b *Bot) handle(ctx context.Context) {
	for {
		select {
		case message := <-b.inbound:
			// persist if/when we need to
			// b.store.SaveMessage(ctx, message)

			b.HandleMessage(message)

		case <-ctx.Done():
			return
		}
	}
}

func (b *Bot) handleText(msg Message) {

}

// todo add a broadcast message for all the users eg updates and what not
