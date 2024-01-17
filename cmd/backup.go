package cmd

import (
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	palSync "github.com/techierishi/pal/sync"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup all configs",
	Long:  `backup - Backup all configs to passed dir`,
	RunE:  backupFunc,
}

func init() {
	RootCmd.AddCommand(backupCmd)
}

func backupFunc(cmd *cobra.Command, args []string) (err error) {

	return palSync.BackupFiles([]string{
		config.Conf.General.SnippetFile,
		config.Conf.General.CredFile,
		config.Conf.General.AliasFile,
	}, config.Conf.General.BackupFile)
}
