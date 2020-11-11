package bot

import (
	"strconv"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

func (b *Bot) isAuthorized(u *tg.Update) bool {
	return u.CallbackQuery != nil &&
		strings.HasPrefix(u.CallbackQuery.Data, CallbackQueryConfirm)
}

func (b *Bot) handleAuthorized(
	l zerolog.Logger,
	u *tg.Update,
) {
	// Parse data
	expectedID, err := strconv.Atoi(
		u.CallbackQuery.Data[len(CallbackQueryConfirm)+1:],
	)
	if err != nil {
		l.Error().Err(err).Msg("parsing user ID from callback query data")
		return
	}

	fromUser := u.CallbackQuery.From

	if fromUser.ID != expectedID {
		// Permission denied, user isn't the one who joined
		l.Debug().
			Int("userClickedID", fromUser.ID).
			Int("userJoinedID", expectedID).
			Msg("permission denied")
		return
	}

	// Authorize user
	tm, ok := b.pending[expectedID]
	if !ok {
		return
	}
	delete(b.pending, expectedID)
	b.stats.UsersAuthorized++
	tm.Stop()

	l.Info().
		Int("userID", fromUser.ID).
		Str("userUsername", fromUser.UserName).
		Msg("authorized user")

	// Delete confirmation request message
	b.logged.DeleteMessage(
		u.CallbackQuery.Message.Chat.ID,
		u.CallbackQuery.Message.MessageID,
		"",
	)
}
