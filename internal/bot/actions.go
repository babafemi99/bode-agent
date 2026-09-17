package bot

func (b *Bot) HandleMessage(msg Message) {
	if msg.IsCommand {
		b.handleCommand(msg)
		return
	}

	b.handleText(msg)
}

func (b *Bot) handleCommand(msg Message) {
	switch msg.Command {
	case "start":
		b.handleStart(msg)

	case "help":
		b.handleHelp(msg)

	case "context":
		b.handleContext(msg)

	default:
		b.handleUnknownCommand(msg)
	}
}
