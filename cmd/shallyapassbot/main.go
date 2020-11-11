package main

import (
	"flag"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"shallyapassbot/bot"
)

func main() {
	flagAdmins := flag.String("admins", "", "bot administrator usernames")
	flagLogFormat := flag.String(
		"logfmt",
		ConfLogFormatConsole,
		"log output format",
	)
	flagDebug := flag.Bool("debug", false, "enable debugging logs")
	flagConfirmTimeout := flag.Duration(
		"confirm-timeout",
		0,
		"confirmation timeout duration. "+
			"Confirmation won't timeout if confirm-timeout is 0",
	)
	flagBanPeriod := flag.Duration(
		"ban-period",
		0,
		"duration for which banned users remain banned. "+
			"Users will be banned indefinitely if ban-period is 0",
	)
	flagTextFmtWelcome := flag.String(
		"text-fmt-welcome",
		"Welcome, [{name}](tg://user?id={id})\\! Please confirm you're human "+
			"by clicking the button below within {confirm-timeout}",
		"welcome message text format. Available placeholders: "+
			"{name} (user's first and last names), "+
			"{id} (users internal Telegram ID), "+
			"{confirm-timeout} (confirmation timeout duration)",
	)
	flagTextConfirm := flag.String(
		"text-confirm",
		"Confirm",
		"confirmation button text",
	)

	flag.Parse()

	l := initLog(
		*flagLogFormat,
		*flagDebug,
	)

	conf := &bot.Config{
		Debug:    *flagDebug,
		APIToken: os.Getenv("TOKEN"),
		AdminUsernames: strings.FieldsFunc(
			*flagAdmins,
			func(r rune) bool { return r == ',' },
		),
		TextFmtWelcome:    *flagTextFmtWelcome,
		TextConfirmButton: *flagTextConfirm,
		ConfirmTimeout:    *flagConfirmTimeout,
		BanPeriod:         *flagBanPeriod,
	}

	b, err := bot.NewBot(conf, l)
	if err != nil {
		l.Fatal().Err(err).Msg("initializing bot API client")
	}

	{
		l.Info().
			Str("username", b.Username()).
			Int("ID", b.ID()).
			Strs("admins", conf.AdminUsernames).
			Bool("conf.debug", conf.Debug).
			Dur("conf.confirmTimeout", conf.ConfirmTimeout).
			Dur("conf.banPeriod", conf.BanPeriod).
			Str("conf.textConfirmButton", conf.TextConfirmButton).
			Str("conf.textFmtWelcome", conf.TextFmtWelcome).
			Msg("initialized")
	}

	if err := b.ListenUpdates(); err != nil {
		l.Fatal().Err(err).Msg("listening for updates")
	}

	l.Info().Msg("shutting down")
}

func initLog(
	format string,
	debug bool,
) zerolog.Logger {
	// Use non-debug console logger by default
	l := zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stdout,
	}).Level(zerolog.InfoLevel)

	switch format {
	case ConfLogFormatJSON:
	case ConfLogFormatConsole:
	case "":
		format = ConfLogFormatConsole
	default:
		l.Fatal().Str("value", format).Msg("invalid log format")
	}

	if format == ConfLogFormatJSON {
		l = zerolog.New(os.Stdout)
	}
	if debug {
		l = l.Level(zerolog.DebugLevel)
	}

	return l.With().Timestamp().Logger()
}

const (
	// ConfLogFormatConsole defines pretty-printed human readable logs
	ConfLogFormatConsole = "console"

	// ConfLogFormatJSON defines pretty raw json logs
	ConfLogFormatJSON = "json"
)
