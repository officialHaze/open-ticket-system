package notification

import (
	"Slack/cmd" // ← now cmd is actually used
	"Slack/notifier"
	"fmt"
	"log"
	"ots/model"
	"ots/settings"
)
type SlackNotifier struct {
	client *notifier.Client
}
var Default *SlackNotifier
func Init(){
	cmd.LoadConfig()
	Default=&SlackNotifier{}
	log.Println("slack notifier initialized via CLI package")
}

func (s *SlackNotifier) NotifyNewTicket(channelID string, t *model.Ticket) {
	if s == nil {
		log.Println("slack notifier not initialized, skipping")
		return
	}
	if channelID == "" {
		log.Println("slack: no channel ID provided, skipping")
		return
	}
	msg := fmt.Sprintf(
		"🎫 *New Ticket Generated*\n*Ticket ID:* `%s`\n*Title:* %s\n*Priority:* %s\n*Status:* %s\n*Created by:* %s",
		t.ID.Hex(), t.Title, t.Priority, t.Status, t.CreatorId,
	)
	if err := cmd.SendMessage(settings.MySettings.Get_SlackBotToken(), channelID, msg); err != nil {
		log.Printf("slack: failed for ticket %s: %v", t.ID.Hex(), err)
	} else {
		log.Printf("slack: notified for ticket %s to channel %s", t.ID.Hex(), channelID)
	}
}