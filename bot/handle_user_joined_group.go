package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

func (b *Bot) isUserJoinedGroup(u *tg.Update) bool {
	return u.Message != nil &&
		u.Message.NewChatMembers != nil &&
		(*u.Message.NewChatMembers)[0].ID != b.api.Self.ID
}

func (b *Bot) handleUserJoinedGroup(
	l zerolog.Logger,
	u *tg.Update,
) {
	user := (*u.Message.NewChatMembers)[0]

	b.stats.UsersJoined++

	if _, ok := b.pending[user.ID]; ok {
		return
	}

	chatID := u.Message.Chat.ID

	msg := tg.NewMessage(
		chatID,
		strings.NewReplacer(
			"{name}", user.String(),
			"{id}", strconv.Itoa(user.ID),
			"{confirm-timeout}", b.conf.ConfirmTimeout.String(),
		).Replace(b.conf.TextFmtWelcome),
	)
	msg.ReplyToMessageID = u.Message.MessageID
	d := fmt.Sprintf(
		"%s_%d",
		CallbackQueryConfirm, user.ID,
	)
	msg.ReplyMarkup = tg.NewInlineKeyboardMarkup(
		[]tg.InlineKeyboardButton{
			{
				Text:         b.conf.TextConfirmButton,
				CallbackData: &d,
			},
		},
	)
	msg.ParseMode = "MarkdownV2"
	msgSent, err := b.logged.SendMessage(
		msg,
		"sending confirmation request message",
	)
	if err != nil {
		return
	}

	timeoutDur := b.conf.ConfirmTimeout
	b.pending[user.ID] = time.AfterFunc(timeoutDur, func() {
		b.lock.Lock()
		defer b.lock.Unlock()

		if _, ok := b.pending[user.ID]; !ok {
			return
		}
		delete(b.pending, user.ID)

		b.log.Info().
			Int64("chatID", chatID).
			Int("userID", user.ID).
			Dur("timeoutDur", timeoutDur).
			Msg("kicking member for exceeding timeout")

		// Kick user for exceeding timeout
		c := tg.KickChatMemberConfig{
			ChatMemberConfig: tg.ChatMemberConfig{
				ChatID: chatID,
				UserID: user.ID,
			},
		}
		if b.conf.BanPeriod != 0 {
			c.UntilDate = time.Now().Add(b.conf.BanPeriod).Unix()
		}
		if _, err := b.logged.KickChatMember(c, ""); err != nil {
			return
		}

		// Remove authorization message
		b.logged.DeleteMessage(
			chatID, msgSent.MessageID,
			"removing authorization message",
		)

		b.stats.UsersBanned++
	})
}
