package bot

import (
	"context"

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
