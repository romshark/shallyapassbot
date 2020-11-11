package bot

import (
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

// logged wraps common API methods with error logs
type logged struct {
	api *tg.BotAPI
	log zerolog.Logger
}

func (l *logged) DeleteMessage(
	chatID int64, messageID int,
	action string,
) (r tg.APIResponse, err error) {
	if action == "" {
		action = "deleting message"
	}
	if r, err = l.api.DeleteMessage(tg.DeleteMessageConfig{
		ChatID:    chatID,
		MessageID: messageID,
	}); err != nil {
		l.log.Error().
			Err(err).
			Int64("chatID", chatID).
			Int("messageID", messageID).
			Msg(action)
	}
	return
}

func (l *logged) KickChatMember(
	c tg.KickChatMemberConfig,
	action string,
) (r tg.APIResponse, err error) {
	if action == "" {
		action = "kicking chat member"
	}
	log := l.log
	if c.UntilDate != 0 {
		tm := time.Unix(c.UntilDate, 0)
		log = log.With().
			Time("untilTime", tm).
			Int64("untilDate", c.UntilDate).
			Logger()
	}
	if r, err = l.api.KickChatMember(c); err != nil {
		log.Error().
			Err(err).
			Int64("chatID", c.ChatID).
			Int("userID", c.UserID).
			Msg(action)
	}
	return
}

func (l *logged) SendMessage(
	c tg.MessageConfig,
	action string,
) (m tg.Message, err error) {
	if action == "" {
		action = "sending message"
	}
	if m, err = l.api.Send(c); err != nil {
		l.log.Error().
			Err(err).
			Int64("chatID", c.ChatID).
			Str("text", c.Text).
			Str("parseMode", c.ParseMode).
			Msg(action)
	}
	return
}
