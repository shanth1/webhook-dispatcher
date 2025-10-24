package config

type Client struct {
	Name   string `mapstructure:"name"`
	ChatID string `mapstructure:"chat_id"`
}

type Bot struct {
	Token   string   `mapstructure:"token"`
	Clients []Client `mapstructure:"clients"`
}

type Telegram struct {
	Bots []Bot `mapstructure:"bots"`
}

type Config struct {
	Addr          string    `mapstructure:"addr"`
	WebhookSecret string    `mapstructure:"webhook_secret"`
	Telegram      *Telegram `mapstructure:"telegram"`
}
