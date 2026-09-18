package bot

import (
	"encoding/json"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WebhookHandler struct {
	bot *Bot
}

func NewWebhookHandler(bot *Bot) *WebhookHandler {
	return &WebhookHandler{
		bot: bot,
	}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var update tgbotapi.Update

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if update.Message == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	msg := Message{
		BotID:     h.bot.ID,
		UserID:    update.Message.From.ID,
		ChatID:    update.Message.Chat.ID,
		Text:      update.Message.Text,
		IsCommand: update.Message.IsCommand(),
		Command:   update.Message.Command(),
		Arguments: update.Message.CommandArguments(),
	}

	select {
	case h.bot.inbound <- msg:
	default:
		http.Error(w, "queue full", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}
