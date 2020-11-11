package bot

import (
	"fmt"
	"time"
)

// Config is the bot's configuration
type Config struct {
	APIToken          string
	AdminUsernames    []string
	Debug             bool
	TextFmtWelcome    string
	TextConfirmButton string
	ConfirmTimeout    time.Duration
	BanPeriod         time.Duration
}

// Copy returns a deep copy of the config
func (c *Config) Copy() *Config {
	cp := *c

	cp.AdminUsernames = make([]string, len(c.AdminUsernames))
	copy(cp.AdminUsernames, c.AdminUsernames)

	return &cp
}

// ValidateAndSetDefaults validates configuration values
// and sets default values if necessary
func (c *Config) ValidateAndSetDefaults() error {
	if c.APIToken == "" {
		return fmt.Errorf("missing access token")
	}

	if c.TextFmtWelcome == "" {
		return fmt.Errorf("empty value for TextFmtWelcome")
	}

	if c.BanPeriod < 30*time.Second {
		c.BanPeriod = 0
	}

	if c.ConfirmTimeout < time.Millisecond {
		return fmt.Errorf(
			"invalid confirm timeout (%s)",
			c.ConfirmTimeout,
		)
	}

	return nil
}

// IsAdmin returns true if the user identified by the given username
// is configured to be one of the bot's administrators
func (c *Config) IsAdmin(username string) bool {
	for _, n := range c.AdminUsernames {
		if n == username {
			return true
		}
	}
	return false
}
