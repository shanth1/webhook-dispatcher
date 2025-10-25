package config

import "github.com/mitchellh/mapstructure"

type SenderType string

const (
	SenderTypeTelegram SenderType = "telegram"
	SenderTypeEmail    SenderType = "email"
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

type Config struct {
	Addr          string         `mapstructure:"addr"`
	WebhookSecret string         `mapstructure:"webhook_secret"`
	Senders       []SenderConfig `mapstructure:"senders"`
	Recipients    []Recipient    `mapstructure:"recipients"`
}

func (sc *SenderConfig) DecodeSenderSettings(v interface{}) error {
	return mapstructure.Decode(sc.Settings, v)
}
