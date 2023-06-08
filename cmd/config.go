package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/autoproof/cli/project"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to detect current directory: %w", err)
		}
		return checkAutoproofProject(currentDir)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to detect current directory: %w", err)
		}

		p, err := project.FromPath(currentDir)
		if err != nil {
			return fmt.Errorf("failed to open project: %w", err)
		}

		config := p.Config()

		switch {
		case len(args) == 0:
			for k, v := range config.AllSettings() {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s = %s\n", k, v)
			}

			return nil

		case len(args) == 1:
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", config.Get(args[0]))
			return nil
		case len(args) == 2:
			config.Set(args[0], args[1])
			return config.Save()
		default:
			return errors.New("invalid arguments number")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
