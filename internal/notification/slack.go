package notification

import (
	"Slack/cmd"
	"Slack/notifier"
	"fmt"
	"log"
	"os"
	"ots/model"
)

type SlackNotifier struct {
	client *notifier.Client
}

var Default *SlackNotifier

func init() {
	cmd.LoadConfig()
	Default = &SlackNotifier{}
	log.Println("slack notifier initialized via CLI package")
}

func (s *SlackNotifier) NotifyNewTicket(channelID string, t *model.Ticket, resolver *model.Resolver) {
	if s == nil {
		log.Println("slack notifier not initialized, skipping")
		return
	}
	if channelID == "" {
		log.Println("slack: no channel ID provided, skipping")
		return
	}

	// Set truncate val
	truncateVal := min(len(t.Description), 100)

	msg := fmt.Sprintf(
		"*New Ticket Generated*\n*Ticket ID:* `%s`\n*Title:* %s\n*Description:* %s\n*Status:* %s\n*Created by:* %s\n*Assigned To:* %s\n*Assignee Email:* %s\n",
		t.ID.Hex(), t.Title, t.Description[:truncateVal], t.Status, t.CreatorId, resolver.Name, resolver.Email,
	)
	if err := cmd.SendMessage(os.Getenv("SLACK_BOT_TOKEN"), channelID, msg); err != nil {
		log.Printf("slack: failed for ticket %s: %v", t.ID.Hex(), err)
	} else {
		log.Printf("slack: notified for ticket %s to channel %s", t.ID.Hex(), channelID)
	}
}
