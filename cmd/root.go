package cmd

import (
	"os"

	"github.com/autoproof/cli/project"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "autoproof",
	Short: `Automatic code & content protection tool.

See 'https://autoproof.dev/docs' for an overview of the system.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func checkAutoproofProject(d string) error {
	_, err := project.FromPath(d)
	return err
}
