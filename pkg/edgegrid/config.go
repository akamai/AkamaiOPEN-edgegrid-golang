package edgegrid

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"
)

const (
	// DefaultConfigFile is the default configuration file path
	DefaultConfigFile = "~/.edgerc"

	// DefaultSection is the .edgerc ini default section
	DefaultSection = "default"

	// MaxBodySize is the max payload size for client requests
	MaxBodySize = 131072
)

type (
	// Config struct provides all the necessary fields to
	// create authorization header, debug is optional
	Config struct {
		Host         string   `ini:"host"`
		ClientToken  string   `ini:"client_token"`
		ClientSecret string   `ini:"client_secret"`
		AccessToken  string   `ini:"access_token"`
		AccountKey   string   `ini:"account_key"`
		HeaderToSign []string `ini:"headers_to_sign"`
		MaxBody      int      `ini:"max_body"`
		Debug        bool     `ini:"debug"`

		file    string
		section string
		env     bool
	}

	// Option defines a configuration option
	Option func(*Config)
)

// New returns new configuration with the specified options
func New(opts ...Option) (*Config, error) {
	c := &Config{
		file:    DefaultConfigFile,
		section: DefaultSection,
		env:     false,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.env {
		if err := c.FromEnv(c.section); err == nil {
			return c, nil
		}
	}

	if err := c.FromFile(c.file, c.section); err != nil {
		return c, fmt.Errorf("Unable to load config from environment or .edgerc file: %w", err)
	}

	return c, nil
}

// WithFile sets the config file path
func WithFile(file string) Option {
	return func(c *Config) {
		c.file = file
	}
}

// WithSection sets the section in the config
func WithSection(section string) Option {
	return func(c *Config) {
		c.section = section
	}
}

// WithEnv sets the option to try to the environment vars to populate the config
// If loading from the env fails, will fallback to .edgerc
func WithEnv(env bool) Option {
	return func(c *Config) {
		c.env = env
	}
}

// FromFile creates a config the configuration in standard INI forma
func (c *Config) FromFile(file string, section string) error {
	var (
		requiredOptions = []string{"host", "client_token", "client_secret", "access_token"}
	)

	path, err := homedir.Expand(file)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	edgerc, err := ini.Load(path)
	if err != nil {
		return fmt.Errorf("could not load config file: %w", err)
	}

	sec, err := edgerc.GetSection(section)
	if err != nil {
		return err
	}

	err = sec.MapTo(&c)
	if err != nil {
		return err
	}

	for _, opt := range requiredOptions {
		if !(edgerc.Section(section).HasKey(opt)) {
			return fmt.Errorf("required option %q is missing from edgerc", opt)
		}
	}

	if c.MaxBody == 0 {
		c.MaxBody = MaxBodySize
	}

	return nil
}

// FromEnv creates a new config using the Environment (ENV)
//
// By default, it uses AKAMAI_HOST, AKAMAI_CLIENT_TOKEN, AKAMAI_CLIENT_SECRET,
// AKAMAI_ACCESS_TOKEN, and AKAMAI_MAX_BODY variables.
//
// You can define multiple configurations by prefixing with the section name specified, e.g.
// passing "ccu" will cause it to look for AKAMAI_CCU_HOST, etc.
//
// If AKAMAI_{SECTION} does not exist, it will fall back to just AKAMAI_.
func (c *Config) FromEnv(section string) error {
	var (
		requiredOptions = []string{"HOST", "CLIENT_TOKEN", "CLIENT_SECRET", "ACCESS_TOKEN"}
		prefix          string
	)

	section = strings.ToUpper(section)

	prefix = "AKAMAI"

	if section != DefaultSection {
		prefix = "AKAMAI_" + strings.ToUpper(section)
	}

	for _, opt := range requiredOptions {
		optKey := fmt.Sprintf("%s_%s", prefix, opt)

		val, ok := os.LookupEnv(optKey)
		if !ok {
			return fmt.Errorf("required option %q is missing from env", optKey)
		}
		switch {
		case opt == "HOST":
			c.Host = val
		case opt == "CLIENT_TOKEN":
			c.ClientToken = val
		case opt == "CLIENT_SECRET":
			c.ClientSecret = val
		case opt == "ACCESS_TOKEN":
			c.AccessToken = val
		}
	}

	c.MaxBody = 0

	val := os.Getenv(fmt.Sprintf("%s_%s", prefix, "MAX_BODY"))
	if i, err := strconv.Atoi(val); err == nil {
		c.MaxBody = i
	}

	if c.MaxBody <= 0 {
		c.MaxBody = MaxBodySize
	}

	return nil
}

// Timestamp returns an edgegrid timestamp from the time
func Timestamp(t time.Time) string {
	local := time.FixedZone("GMT", 0)
	t = t.In(local)
	return fmt.Sprintf("%d%02d%02dT%02d:%02d:%02d+0000",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}