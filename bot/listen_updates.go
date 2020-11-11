package bot

import (
	"fmt"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

// ListenUpdates starts listening for updates
// blocking until StopListening is called
func (b *Bot) ListenUpdates() error {
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("getting update channel: %w", err)
	}

LISTEN:
	for u := range updates {
		l := b.log.With().
			Int("updateID", u.UpdateID).
			Time("updateTime", time.Now()).
			Logger()

		byAdmin := false
		if u.Message != nil {
			byAdmin = b.conf.IsAdmin(u.Message.From.UserName)
		}

		l.Debug().
			Bool("byAdmin", byAdmin).
			Interface("update", u).
			Msg("received update")

		for _, h := range b.handles {
			l := l.With().Str("updateType", h.Name).Logger()
			if func() (processed bool) {
				b.lock.Lock()
				defer b.lock.Unlock()

				if !h.Select(&u) {
					return false
				}
				l.Info().Msg("handling update")
				h.Process(l, &u)
				return true
			}() {
				continue LISTEN
			}
		}
		l.Debug().Msg("no update handler selected")
	}

	return nil
}
