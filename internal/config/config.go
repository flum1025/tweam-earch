package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Accounts Accounts `yaml:"accounts"`
}

type Account struct {
	ID    string `yaml:"id"`
	Slack Slack  `yaml:"slack"`
	Rules Rules  `yaml:"rules"`
}

type Accounts []Account

func (a Accounts) Find(userID string) *Account {
	for _, account := range a {
		if account.ID == userID {
			return &account
		}
	}

	return nil
}

type Slack struct {
	Token      string `yaml:"token"`
	UserIcon   string `yaml:"user_icon"`
	Channel    string `yaml:"channel"`
	Icon       string `yaml:"icon"`
	SourceUser string `yaml:"source_user"`
}

type Rule struct {
	Text string `yaml:"text"`
}

type Rules []Rule

func NewConfig(configPath string) (*Config, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	config := Config{}

	if err = yaml.Unmarshal(buf, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if len(config.Accounts) == 0 {
		return nil, fmt.Errorf("accounts are required")
	}

	m := make(map[string]Account)

	for _, account := range config.Accounts {
		if _, ok := m[account.ID]; ok {
			return nil, fmt.Errorf("deplicate account IDs")
		}

		m[account.ID] = account
	}

	return &config, nil
}
