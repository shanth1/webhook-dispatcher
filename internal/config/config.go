package config

type Client struct {
	Name   string `yaml:"name"`
	ChatID string `yaml:"chat_id"`
}

type Bot struct {
	Token   string   `yaml:"token"`
	Clients []Client `yaml:"clients"`
}

type Telegram struct {
	Bots []Bot `yaml:"bots"`
}

type Config struct {
	Addr          string    `yaml:"addr"`
	WebhookSecret string    `yaml:"webhook_secret"`
	Telegram      *Telegram `yaml:"telegram"`
}
