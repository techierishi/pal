package cmd

import (
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	palSync "github.com/techierishi/pal/sync"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync configs",
	Long:  `Sync configs with gist/gitlab`,
	RunE:  syncFunc,
}

func init() {
	RootCmd.AddCommand(syncCmd)
}

func syncFunc(cmd *cobra.Command, args []string) (err error) {
	err = palSync.BackupFiles([]string{
		config.Conf.General.SnippetFile,
		config.Conf.General.CredFile,
		config.Conf.General.AliasFile,
	}, config.Conf.General.BackupFile)

	if err != nil {
		return
	}

	return palSync.AutoSync(config.Conf.General.BackupFile)
}
