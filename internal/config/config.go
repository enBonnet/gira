package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	JiraURL  string `mapstructure:"jira_url"`
	Email    string `mapstructure:"email"`
	APIToken string `mapstructure:"api_token"`
	Project  string `mapstructure:"project"`
}

func Load() (*Config, error) {
	viper.SetConfigName("gira")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/gira")
	viper.AddConfigPath("$HOME/.gira")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("GIRA")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
