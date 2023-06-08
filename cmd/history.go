package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Displays a message about retrieving the history of your project",
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to detect current directory: %w", err)
		}
		return checkAutoproofProject(currentDir)
	},
	Run: func(cmd *cobra.Command, _ []string) {
		_, _ = fmt.Fprintln(
			cmd.OutOrStdout(),
			"See all history of your project <https://app.autoproof.dev/history>.",
		)
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)
}
