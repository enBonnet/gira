package cmd

import (
	"fmt"

	"gira/internal/config"
	"gira/internal/jira"

	"github.com/spf13/cobra"
)

var myselfCmd = &cobra.Command{
	Use:   "myself",
	Short: "Show current user info",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.JiraURL == "" || cfg.Email == "" || cfg.APIToken == "" {
			return fmt.Errorf("missing configuration")
		}

		client := jira.NewClient(cfg.JiraURL, cfg.Email, cfg.APIToken)

		user, err := client.GetMyself()
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		fmt.Printf("Account ID: %s\n", user.AccountID)
		fmt.Printf("Name:       %s\n", user.DisplayName)
		fmt.Printf("Email:      %s\n", user.EmailAddress)
		fmt.Printf("Active:     %v\n", user.Active)
		fmt.Printf("Timezone:   %s\n", user.TimeZone)
		fmt.Printf("Locale:     %s\n", user.Locale)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(myselfCmd)
}
