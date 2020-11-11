package bot

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

func (b *Bot) isAddedToGroup(u *tg.Update) bool {
	return u.Message != nil &&
		u.Message.NewChatMembers != nil &&
		(*u.Message.NewChatMembers)[0].ID == b.api.Self.ID
}

func (b *Bot) handleAddedToGroup(
	l zerolog.Logger,
	u *tg.Update,
) {
	// TODO
}
