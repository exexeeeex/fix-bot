package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go/v2"
	"github.com/max-messenger/max-bot-api-client-go/v2/model"
)

func bot(ctx context.Context, token string) {
	token = strings.TrimSpace(token)

	if token == "" {
		log.Println("[MAX] MAX_BOT_TOKEN is empty")
		return
	}

	opts := []maxbot.Opt{
		maxbot.WithHTTPClient(&http.Client{
			Timeout: 60 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Только для DEV!
				},
			},
		}),
	}

	api, err := maxbot.NewApi(token, opts...)
	if err != nil {
		log.Printf("[MAX] NewApi error: %v", err)
		return
	}

	info, err := api.Bots.GetMyInfo(ctx)
	if err != nil {
		log.Printf("[MAX] GetMyInfo error: %v", err)
		return
	}

	log.Printf("[MAX] Bot authorized: %+v", info)

	subscriptions, err := api.Subscriptions.GetSubscriptions(ctx)
	if err != nil {
		log.Printf("[MAX] GetSubscriptions error: %v", err)
		return
	}

	if len(subscriptions.Subscriptions) > 0 {
		log.Println("[MAX] Webhook subscription exists.")
		log.Println("[MAX] Long Polling cannot be used simultaneously with Webhook.")
		return
	}

	log.Println("[MAX] Starting Long Polling...")

	handle := func(ctx context.Context, update model.Update) {
		log.Printf(
			"[MAX] Received update: type=%s chat=%d user=%d",
			update.UpdateType,
			update.ChatID,
			update.UserID,
		)

		switch update.UpdateType {
		case model.UpdateMessageCreated:
			// cmd := update.GetCommand()
			keyboard := model.NewKeyboard()
			keyboard.AddRow().AddOpenApp("🚀 Открыть Fiks", info.UserID)

			msg := maxbot.NewMessage().
				AddKeyboard(keyboard).
				SetText("Привет! 👋\nЯ бот приложения Fiks. Чем могу помочь?").
				SetChat(update.ChatID).
				SetUser(update.UserID)

			result, err := api.Messages.Send(ctx, msg)
			if err != nil {
				log.Printf("[MAX] Send message error: %v", err)
				return
			}

			log.Printf("[MAX] Message sent: %+v", result)

		default:
			log.Printf("[MAX] Unsupported update type: %s", update.UpdateType)
		}
	}

	var marker int64

	for {
		if ctx.Err() != nil {
			log.Println("[MAX] Bot stopped")
			return
		}

		updates, nextMarker, err := api.Subscriptions.GetUpdates(ctx, marker)

		if err != nil {
			if ctx.Err() != nil {
				log.Println("[MAX] Bot stopped")
				return
			}

			var timeoutErr *maxbot.TimeoutError
			if errors.As(err, &timeoutErr) {
				continue
			}

			log.Printf("[MAX] GetUpdates error: %v", err)

			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
				continue
			}
		}

		marker = nextMarker

		for _, update := range updates {
			if ctx.Err() != nil {
				return
			}

			handle(ctx, update)
		}
	}
}
