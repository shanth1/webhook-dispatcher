package config

import "github.com/mitchellh/mapstructure"

type SenderType string

const (
	SenderTypeTelegram SenderType = "telegram"
	SenderTypeEmail    SenderType = "email"
)

type WebhookType string

const (
	WebhookTypeGitHub   WebhookType = "github"
	WebhookTypeKanboard WebhookType = "kanboard"
	WebhookTypeCustom   WebhookType = "custom"
)

type Recipient struct {
	Name   string `mapstructure:"name"`
	Sender string `mapstructure:"sender"`
	Target string `mapstructure:"target"`
}

type SenderConfig struct {
	Name     string                 `mapstructure:"name"`
	Type     SenderType             `mapstructure:"type"`
	Settings map[string]interface{} `mapstructure:"settings"`
}

type WebhookConfig struct {
	Name       string      `mapstructure:"name"`
	Path       string      `mapstructure:"path"`
	Type       WebhookType `mapstructure:"type"`
	Secret     string      `mapstructure:"secret"`
	Recipients []string    `mapstructure:"recipients"`
}

type TelegramSettings struct {
	Token string `mapstructure:"token"`
}

type EmailSettings struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

type Logger struct {
	App       string `mapstructure:"app"`
	Level     string `mapstructure:"level"`
	Service   string `mapstructure:"service"`
	UDPAddres string `mapstructure:"udp_address"`
}

type Config struct {
	Env        string          `mapstructure:"env"`
	Addr       string          `mapstructure:"addr"`
	Webhooks   []WebhookConfig `mapstructure:"webhooks"`
	Senders    []SenderConfig  `mapstructure:"senders"`
	Recipients []Recipient     `mapstructure:"recipients"`
	Logger     Logger          `mapstructure:"logger"`
}

func (sc *SenderConfig) DecodeSenderSettings(v interface{}) error {
	return mapstructure.Decode(sc.Settings, v)
}
