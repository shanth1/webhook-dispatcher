package config

import "github.com/mitchellh/mapstructure"

type NotifierType string
type WebhookType string

type NotifierName string
type WebhookName string

const (
	NotifierTypeTelegram NotifierType = "telegram"
	NotifierTypeEmail    NotifierType = "email"
)

const (
	WebhookTypeGitHub   WebhookType = "github"
	WebhookTypeKanboard WebhookType = "kanboard"
	WebhookTypeCustom   WebhookType = "custom"
)

type Config struct {
	Env                     string           `mapstructure:"env"`
	Addr                    string           `mapstructure:"addr"`
	Webhooks                []WebhookConfig  `mapstructure:"webhooks"`
	Notifiers               []NotifierConfig `mapstructure:"notifiers"`
	Recipients              []Recipient      `mapstructure:"recipients"`
	Logger                  Logger           `mapstructure:"logger"`
	DisableUnknownTemplates bool             `mapstructure:"disable_unknown_templates"` // TODO: moved to webhookConfig
}

type Logger struct {
	App        string `mapstructure:"app"`
	Level      string `mapstructure:"level"`
	Service    string `mapstructure:"service"`
	UDPAddress string `mapstructure:"udp_address"`
}

type NotifierConfig struct {
	Name     NotifierName           `mapstructure:"name"`
	Type     NotifierType           `mapstructure:"type"`
	Settings map[string]interface{} `mapstructure:"settings"`
}

func (sc *NotifierConfig) DecodeSettings(v interface{}) error {
	return mapstructure.Decode(sc.Settings, v)
}

type WebhookConfig struct {
	Name       WebhookName `mapstructure:"name"`
	Path       string      `mapstructure:"path"`
	Type       WebhookType `mapstructure:"type"`
	Secret     string      `mapstructure:"secret"`
	BaseURL    string      `mapstructure:"base_url"`
	Recipients []string    `mapstructure:"recipients"`
}

type Recipient struct {
	Name     string       `mapstructure:"name"`
	Target   string       `mapstructure:"target"`
	Notifier NotifierName `mapstructure:"notifier"`
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
