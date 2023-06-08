package project

import (
	"errors"
	"path"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

func (p *Project) Config() *Config {
	return p.config
}

func (p *Project) readConfig() error {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("AUTOPROOF")
	v.SetConfigFile(path.Join(p.path, AutoproofHomeDirName, "config.yml"))

	p.config = &Config{
		viper: v,
	}

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func (p *Project) writeConfig() error {
	if p.config == nil || p.config.viper == nil {
		return errors.New("incorrect initialization of configuration")
	}

	return nil
}

func (c *Config) AllSettings() map[string]any {
	return c.viper.AllSettings()
}

func (c *Config) Set(key string, value any) {
	c.viper.Set(key, value)
}

func (c *Config) Get(key string) any {
	return c.viper.Get(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

func (c *Config) Save() error {
	c.viper.ConfigFileUsed()
	return c.viper.WriteConfig()
}
