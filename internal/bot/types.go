package bot

type Message struct {
	BotID     string
	UserID    int64
	ChatID    int64
	Text      string
	IsCommand bool
	Command   string
	Arguments string
}
