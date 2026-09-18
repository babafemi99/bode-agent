package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/babafemi99/bode-agent/internal/bot"
	"github.com/go-chi/chi/v5"
)

type API struct {
	Port           int64
	Server         *http.Server
	Bot            *bot.Bot
	WebhookHandler *bot.WebhookHandler
}

func (a *API) Serve() error {
	a.Server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.Port),
		ReadTimeout:    5 * time.Second,
		Handler:        a.setUpServerHandler(),
		MaxHeaderBytes: 1024 * 1024,
	}

	return a.Server.ListenAndServe()
}

func (a *API) Shutdown(ctx context.Context) error {
	if a.Server == nil {
		return nil
	}

	return a.Server.Shutdown(ctx)
}

func (a *API) setUpServerHandler() http.Handler {
	r := chi.NewRouter()

	r.Post("/webhook", a.WebhookHandler.ServeHTTP)

	return r
}
