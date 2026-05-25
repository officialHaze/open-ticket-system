package cmd
import (
	"Slack/config"
	"Slack/notifier"
	"fmt"

	"github.com/spf13/cobra"
)
var ChannelId string
var Message string
var Send = &cobra.Command{
	Use:   "send",
	Short: "Ticket generated...",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ChannelId == "" {
			ChannelId = cfg.Slack.ChannelID
		}
		client := notifier.New(cfg)
		if err := client.Send(ChannelId, Message); err != nil {
			return err
		}
		fmt.Printf("Message sent to %s\n", ChannelId)
		return nil
	},
}
func SendMessage(botToken string, channelID string, message string) error {
	client := notifier.New(&config.Config{
		Slack: config.SlackConfig{
			BotToken: botToken,
		},
	})
	return client.Send(channelID, message)
}
func init() {
	Send.Flags().StringVarP(&ChannelId, "channel", "c", "", "channel ID")
	Send.Flags().StringVarP(&Message, "message", "m", "", "Message sent")
	Send.MarkFlagRequired("message")
	rootCmd.AddCommand(Send)
}