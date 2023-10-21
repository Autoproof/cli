package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/viper"
)

var (
	// secretOptionNames is a list of options that will be stored in a secret file.
	secretOptionNames = []string{"apiKey"}
)

func isSecretOptionName(name string) bool {
	return slices.ContainsFunc(secretOptionNames, func(secretOptionName string) bool {
		return strings.EqualFold(name, secretOptionName)
	})
}

// Config is a configuration of the project.
type Config struct {
	// public is a public configuration of the project.
	public *viper.Viper

	// secret is a secret configuration of the project.
	secret *viper.Viper
}

func (p *Project) Config() *Config {
	return p.config
}

// AllSettings returns all settings of the project.
func (c *Config) AllSettings() map[string]any {
	settings := c.public.AllSettings()
	for key, value := range c.secret.AllSettings() {
		// secret settings overrides public settings.
		settings[key] = value
	}
	return settings
}

// Set sets a value for the given configuration key.
func (c *Config) Set(key string, value any) {
	if isSecretOptionName(key) {
		c.secret.Set(key, value)
		return
	}

	c.public.Set(key, value)
}

// Get returns a value for the given configuration key.
func (c *Config) Get(key string) any {
	if isSecretOptionName(key) {
		return c.secret.Get(key)
	}

	return c.public.Get(key)
}

// GetStringSlice returns a string slice for the given configuration key.
func (c *Config) GetStringSlice(key string) []string {
	if isSecretOptionName(key) {
		return c.secret.GetStringSlice(key)
	}

	return c.public.GetStringSlice(key)
}

// Save saves the configuration to the file system.
func (c *Config) Save() error {
	if err := c.secret.WriteConfig(); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("write secret config: %w", err)
		}

		secretDir := filepath.Dir(c.secret.ConfigFileUsed())
		if err := os.MkdirAll(secretDir, 0700); err != nil {
			return fmt.Errorf("create secret config dir: %w", err)
		}

		if err := c.secret.WriteConfig(); err != nil {
			return fmt.Errorf("write secret config: %w", err)
		}
	}

	v := viper.New()
	v.SetConfigFile(c.public.ConfigFileUsed())
	for _, publicConfigKey := range c.public.AllKeys() {
		if isSecretOptionName(publicConfigKey) {
			continue
		}
		v.Set(publicConfigKey, c.public.Get(publicConfigKey))
	}

	return v.WriteConfig()
}

func (p *Project) readConfig() error {
	public := viper.New()
	public.AutomaticEnv()
	public.SetEnvPrefix("AUTOPROOF")
	public.SetConfigFile(path.Join(p.path, AutoproofHomeDirName, "config.yml"))

	if err := public.ReadInConfig(); err != nil {
		return fmt.Errorf("read public config: %w", err)
	}

	projectName := public.GetString("projectname")
	if projectName == "" {
		// Case when project is in initialization stage. At initialization stage project name is known after
		// configuration file is created.
		projectName = p.name
		if projectName == "" {
			return errors.New("project is not initialized correctly. " +
				"Execute `autoproofcli init` command before using `autoproofcli`")
		}
	}

	// secret configuration file is stored in the user's home directory.
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("current user: %w", err)
	}

	secret := viper.New()
	secret.AutomaticEnv()
	secret.SetEnvPrefix("AUTOPROOF")
	secret.SetConfigFile(path.Join(currentUser.HomeDir, AutoproofHomeDirName, projectName, "secret.yml"))

	p.config = &Config{
		public: public,
		secret: secret,
	}

	// secret configuration file is optional.
	if err := secret.ReadInConfig(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("read secret config: %w", err)
	}

	return nil
}
