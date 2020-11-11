package bot

import (
	"fmt"
	"sync"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog"
)

type Bot struct {
	api     *tg.BotAPI
	log     zerolog.Logger
	logged  *logged
	handles [5]handle

	lock    sync.Mutex
	conf    *Config
	pending map[userID]*time.Timer
	stats   *Statistics
}

type userID = int

// NewBot creates a new instance of the shallyapass bot
func NewBot(
	conf *Config,
	log zerolog.Logger,
) (*Bot, error) {
	if err := conf.ValidateAndSetDefaults(); err != nil {
		return nil, fmt.Errorf("preparing config: %w", err)
	}

	api, err := tg.NewBotAPI(conf.APIToken)
	if err != nil {
		return nil, err
	}

	api.Debug = conf.Debug

	b := &Bot{
		api:  api,
		conf: conf,
		log:  log,
		logged: &logged{
			api: api,
			log: log,
		},
		stats:   &Statistics{},
		pending: map[int]*time.Timer{},
	}
	b.handles = [5]handle{
		{"directMessage", b.isDirectMessage, b.handleDirectMessage},
		{"messageSent", b.isMessageSent, b.handleMessageSent},
		{"addedToGroup", b.isAddedToGroup, b.handleAddedToGroup},
		{"userJoinedGroup", b.isUserJoinedGroup, b.handleUserJoinedGroup},
		{"authorized", b.isAuthorized, b.handleAuthorized},
	}

	return b, nil
}

// StopListening stops all update listeners
func (b *Bot) StopListening() {
	b.api.StopReceivingUpdates()
}

// Username returns the bot's username
func (b *Bot) Username() string {
	return b.api.Self.UserName
}

// ID returns the bot's ID
func (b *Bot) ID() int {
	return b.api.Self.ID
}

// Conf returns the bot's configuration
func (b *Bot) Conf() *Config {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.conf.Copy()
}

// Statistics returns the bot's statistics
func (b *Bot) Statistics() *Statistics {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.stats.Copy()
}

type handle struct {
	Name    string
	Select  func(*tg.Update) bool
	Process func(zerolog.Logger, *tg.Update)
}

const CallbackQueryConfirm = "confirm"

type Statistics struct {
	UsersJoined       uint64
	UsersAuthorized   uint64
	UsersBanned       uint64
	MessagesDeleted   uint64
	MessagesProcessed uint64
}

// Copy returns a deep copy of the statistics
func (s *Statistics) Copy() *Statistics {
	c := *s
	return &c
}
