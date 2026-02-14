package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"gira/internal/config"
	"gira/internal/jira"

	"github.com/spf13/cobra"
)

var taskLimit int

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List your Jira tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.JiraURL == "" || cfg.Email == "" || cfg.APIToken == "" {
			return fmt.Errorf("missing configuration. Set JIRA_URL, JIRA_EMAIL, and JIRA_API_TOKEN in config file or environment variables")
		}

		client := jira.NewClient(cfg.JiraURL, cfg.Email, cfg.APIToken)

		issues, err := client.GetMyTasks(cfg.Project)
		if err != nil {
			return fmt.Errorf("failed to fetch tasks: %w", err)
		}

		if len(issues) == 0 {
			fmt.Println("No tasks found")
			return nil
		}

		sort.Slice(issues, func(i, j int) bool {
			return issues[i].Fields.Priority.ID < issues[j].Fields.Priority.ID
		})

		if taskLimit > 0 && len(issues) > taskLimit {
			issues = issues[:taskLimit]
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "KEY\tSTATUS\tPRIORITY\tSUMMARY\tURL")
		for _, issue := range issues {
			url := client.GetIssueURL(issue.Key)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", issue.Key, issue.Fields.Status.Name, issue.Fields.Priority.Name, issue.Fields.Summary, url)
		}
		w.Flush()

		return nil
	},
}

func init() {
	tasksCmd.Flags().IntVarP(&taskLimit, "limit", "l", 5, "maximum number of tasks to display")
	rootCmd.AddCommand(tasksCmd)
}
