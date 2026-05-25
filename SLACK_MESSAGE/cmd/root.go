package cmd
import (
	"Slack/config"
	"log"
	"github.com/spf13/cobra"
)
var cfg *config.Config
var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "Ticket raise",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		LoadConfig()
	},
}
func LoadConfig() {
	var err error
	cfg, err = config.Load("config.yml")
	if err != nil {
		log.Fatal("error loading config:", err)
	}
}
func Execute() {
	rootCmd.Execute()
}