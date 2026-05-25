package cmd

import (
	"Slack/config"
	"log"

	"github.com/spf13/cobra"
)

var cfg *config.Config
var rootCmd = &cobra.Command{
	Use:   "slackalert",
	Short: "Ticket raise alert",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		LoadConfig()
	},
}

func LoadConfig() {
	var err error
	cfg, err = config.Load("slack_config.yml")
	if err != nil {
		log.Fatal("error loading config:", err)
	}
}
func Execute() {
	rootCmd.Execute()
}
