package bot

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

func (b *Bot) isMessageSent(u *tg.Update) bool {
	return u.Message != nil &&
		u.Message.Chat.Type != "private" &&
		u.Message.Text != ""
}

func (b *Bot) handleMessageSent(
	l zerolog.Logger,
	u *tg.Update,
) {
	b.stats.MessagesProcessed++

	if _, ok := b.pending[u.Message.From.ID]; ok {
		b.logged.DeleteMessage(
			u.Message.Chat.ID, u.Message.MessageID,
			"deleting message from pending user",
		)
		b.stats.MessagesDeleted++
	}
}
