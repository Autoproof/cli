package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/autoproof/cli/autoproofapi"
	"github.com/autoproof/cli/project"
)

const (
	dryRunFlagName  = "dry-run"
	messageFlagName = "message"
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: `Create a snapshot of only hashes of files in a project and send it to a copyright registration center`,
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

		apiKey, ok := p.Config().Get("apiKey").(string)
		if !ok {
			return errors.New("failed to read api key: key not found")
		}

		snapshot, err := p.Snapshot(cmd.Context())
		if err != nil {
			return fmt.Errorf("failed to create snapshot: %w", err)
		}

		apiKeyTransport := &autoproofapi.APIKeyTransport{APIKey: apiKey}
		apiClient := autoproofapi.NewClient(autoproofapi.WithCustomClient(apiKeyTransport.Client()))

		mode := autoproofapi.ProductionSnapshotMode
		if dryRun, _ := cmd.Flags().GetBool(dryRunFlagName); dryRun {
			mode = autoproofapi.TestingSnapshotMode
		}

		description, _ := cmd.Flags().GetString(messageFlagName)
		autoproofSnapshot := &autoproofapi.Snapshot{
			Project: p.Name(),

			Description: description,
			Mode:        mode,
		}
		for _, snapshotItem := range snapshot.Items {
			autoproofSnapshot.Data = append(autoproofSnapshot.Data, autoproofapi.SnapshotItem{
				Filename: snapshotItem.Filename,
				Hash:     snapshotItem.Hash,
			})
		}

		resp, err := apiClient.UploadSnapshot(cmd.Context(), autoproofSnapshot)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), resp.Message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)

	snapshotCmd.Flags().StringP(messageFlagName, "m", "", "Short description sent along with the snapshot to the copyright registration center.")
	snapshotCmd.Flags().BoolP(dryRunFlagName, "n", false, "Perform hash saving in either testing or production mode.")
}
