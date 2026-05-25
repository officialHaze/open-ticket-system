package config
import (
	"os"

	"gopkg.in/yaml.v3"
)
type SlackConfig struct {
	BotToken  string `yaml:"bot_token"`
	AppToken  string `yaml:"app_token"`
	ChannelID string `yaml:"ChannelID"`
}
type Config struct {
	Slack SlackConfig `yaml:"slack"`
}
func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}