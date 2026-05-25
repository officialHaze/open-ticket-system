package notifier

import (
	"Slack/config"
	"fmt"

	"github.com/slack-go/slack"
)

type Client struct {
	api *slack.Client
}
func New(cfg *config.Config)*Client{
	return &Client{
		api:slack.New(cfg.Slack.BotToken),
	}
}
func (c *Client) Send(channelID string,message string) error{
	_,_,err:=c.api.PostMessage(
		channelID,
		slack.MsgOptionText(message,false),
	)
	if err!=nil{
		return fmt.Errorf("slack send :%w",err)
	}
	return nil
}