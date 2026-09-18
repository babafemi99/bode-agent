package bot

import "fmt"

func (b *Bot) handleStart(msg Message) {
	b.outbound <- Message{
		ChatID: msg.ChatID,
		Text: `Welcome to Bọ̀dé!

I can help you manage your events, guests, tables, and invitations.

Try:
• /events - View your events
• /context - See your current context
• /help - See what I can do

Or just send me a message.`,
	}
}

func (b *Bot) handleUnknownCommand(msg Message) {
	b.outbound <- Message{
		ChatID: msg.ChatID,
		Text: fmt.Sprintf(
			`I don't recognize /%s yet.

Try one of these:
• /start - Start using Bọ̀dé
• /help - See what I can do
• /events - View your events
• /context - Show your current context

You can also just send me a message.`,
			msg.Command,
		),
	}
}

func (b *Bot) handleText(msg Message) {

}
