package bot

import (
	"context"
	"log/slog"
	"sync"

	"github.com/babafemi99/bode-agent/pkg/lid"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Manager struct {
	mu sync.RWMutex

	bots map[string]*Bot

	register   chan *Bot
	unregister chan string
}

func (m *Manager) HandleNewBot(token string) error {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
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

	m.register <- bot

	go bot.receive(context.Background())
	go bot.send(context.Background())

	return nil
}

func NewManager() *Manager {
	return &Manager{
		bots:       make(map[string]*Bot),
		register:   make(chan *Bot, 256),
		unregister: make(chan string, 256),
	}
}

func (m *Manager) Run(ctx context.Context) {
	for {
		select {
		case bot := <-m.register:
			m.RegisterBot(bot)

		case botID := <-m.unregister:
			m.DeregisterBot(botID)

		case <-ctx.Done():
			return
		}
	}
}

func (m *Manager) Broadcast(message Message) {
	for _, bot := range m.GetBots() {
		bot.outbound <- message
	}
}

func (m *Manager) RegisterBot(bot *Bot) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if oldBot, exists := m.bots[bot.ID]; exists {
		oldBot.Client.StopReceivingUpdates()

		slog.Info(
			"replaced existing bot",
			"bot_id", bot.ID,
		)
	}

	m.bots[bot.ID] = bot

	slog.Info(
		"bot registered",
		"bot_id", bot.ID,
		"telegram_bot_id", bot.BotID,
	)
}

func (m *Manager) DeregisterBot(botID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.bots[botID]
	if !exists {
		slog.Info(
			"bot not found",
			"bot_id", botID,
		)
		return
	}

	delete(m.bots, botID)

	slog.Info(
		"bot deregistered",
		"bot_id", botID,
	)
}

func (m *Manager) GetBot(botID string) (*Bot, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	bot, exists := m.bots[botID]
	return bot, exists
}

func (m *Manager) GetBots() []*Bot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	bots := make([]*Bot, 0, len(m.bots))

	for _, bot := range m.bots {
		bots = append(bots, bot)
	}

	return bots
}
