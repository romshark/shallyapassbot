package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

func (b *Bot) isDirectMessage(u *tg.Update) bool {
	return u.Message != nil &&
		u.Message.Chat.Type == "private"
}

func (b *Bot) handleDirectMessage(
	l zerolog.Logger,
	u *tg.Update,
) {
	if !b.conf.IsAdmin(u.Message.Chat.UserName) {
		// Reply to unauthorized requests
		msg := tg.NewMessage(u.Message.Chat.ID, "unauthorized request")
		msg.ReplyToMessageID = u.Message.MessageID
		b.logged.SendMessage(msg, "replying to unauthorized request")
		return
	}

	l.Info().
		Str("text", u.Message.Text).
		Int("userID", u.Message.From.ID).
		Msg("handling direct message")

	reply := func(
		text string,
		action string,
	) {
		msg := tg.NewMessage(u.Message.Chat.ID, text)
		msg.ReplyToMessageID = u.Message.MessageID
		msg.ParseMode = "MarkdownV2"
		b.logged.SendMessage(msg, action)
	}

	f := strings.Fields(u.Message.Text)
	switch f[0] {
	case "/help":
		reply(
			MsgTxtHelp,
			"replying to command: help",
		)

	case "/stats":
		itoa := func(i uint64) string { return strconv.FormatUint(i, 10) }
		reply(
			strings.NewReplacer(
				"{UsersJoined}", itoa(b.stats.UsersJoined),
				"{UsersAuthorized}", itoa(b.stats.UsersAuthorized),
				"{UsersBanned}", itoa(b.stats.UsersBanned),
				"{MessagesProcessed}", itoa(b.stats.MessagesProcessed),
				"{MessagesDeleted}", itoa(b.stats.MessagesDeleted),
			).Replace(MsgFmtStats),
			"replying to command: /stats",
		)

	case "/setconf":
		if len(f) < 2 {
			reply(
				"Specify the config key",
				"replying to incomplete set command",
			)
			return
		} else if len(f) < 3 {
			reply(
				"Specify the value",
				"replying to incomplete set command",
			)
			return
		}
		switch f[1] {
		case "confirm-timeout":
			d, err := time.ParseDuration(f[2])
			if err != nil {
				reply(
					fmt.Sprintf("parsing duration: %s", err),
					"replying to invalid duration value",
				)
			}
			b.conf.ConfirmTimeout = d
			reply(
				fmt.Sprintf("confirm\\-timeout set to %s", d),
				"replying to setconf confirm-timeout",
			)

		case "ban-period":
			d, err := time.ParseDuration(f[2])
			if err != nil {
				reply(
					fmt.Sprintf("parsing duration: %s", err),
					"replying to invalid duration value",
				)
			}
			b.conf.BanPeriod = d
			reply(
				fmt.Sprintf("ban\\-period set to %s", d),
				"replying to setconf ban-period",
			)

		default:
			reply(
				"Unknown key\n"+
					"\\- confirm\\-timeout \\<duration\\>\n"+
					"\\- ban\\-period \\<duration\\>\n",
				"replying to invalid conf key",
			)
		}

	default:
		reply(
			"Unknown command\n\nHelp:\n"+MsgTxtHelp,
			"replying to unknown command",
		)
	}
}

const MsgTxtHelp = `/help \- display help
/stats \- display statistics
/setconf \<key\> \<value\> \- sets a configuration value by key`

const MsgFmtStats = "users joined: {UsersJoined}\n" +
	"users authorized: {UsersAuthorized}\n" +
	"users banned: {UsersBanned}\n" +
	"messages processed: {MessagesProcessed}\n" +
	"messages deleted: {MessagesDeleted}\n"
