package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/autoproof/cli/project"
	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
)

const (
	projectNameFlagName = "project-name"
	apiKeyFlagName      = "api-key"
)

var (
	defaultIgnoreList = []string{
		".autoproof",
		".DS_Store",
		".git",
		".svn",
		".cvs",
	}
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a project area in a current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}

		if err := checkAutoproofProject(currentDir); err == nil {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Autoproof project already initialized.")
			return nil
		}

		projectName, err := cmd.Flags().GetString(projectNameFlagName)
		if err != nil {
			return err
		}

		if projectName == "" {
			prompt := promptui.Prompt{
				Label:   "Project Name",
				Default: path.Base(currentDir),
			}

			projectName, err = prompt.Run()
			if err != nil {
				return err
			}
		}

		apiKey, err := cmd.Flags().GetString(apiKeyFlagName)
		if err != nil {
			return err
		}

		if apiKey == "" {
			prompt := promptui.Prompt{
				Label: "API key",
				Validate: func(apiKey string) error {
					if strings.TrimSpace(apiKey) == "" {
						return errors.New("API key is missing and required")
					}
					return nil
				},
			}

			apiKey, err = prompt.Run()
			if err != nil {
				return err
			}
		}

		p, err := project.New(projectName, currentDir)
		if err != nil {
			return fmt.Errorf("failed to create new project: %w", err)
		}

		config := p.Config()
		config.Set("apiKey", apiKey)
		config.Set("ignore", defaultIgnoreList)

		if err := config.Save(); err != nil {
			return err
		}

		//if err := os.MkdirAll(path.Join(currentDir, ".autoproof"), 0770); err != nil {
		//	return err
		//}
		//
		//v := viper.New()
		//v.Set("projectName", projectName)
		//v.Set("apiKey", apiKey)
		//v.Set("ignore", defaultIgnoreList)
		//
		//v.AddConfigPath(path.Join(currentDir, ".autoproof"))
		//if err := v.WriteConfigAs(path.Join(currentDir, ".autoproof", "config.yml")); err != nil {
		//	return err
		//}
		//
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Project %q has been successfully initialized.\n", p.Name())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP(apiKeyFlagName, "k", "", "The API key for authentication with the copyright center.")
	initCmd.Flags().StringP(projectNameFlagName, "p", "", "The project name")
}
