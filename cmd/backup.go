package cmd

import (
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/syncm"
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

	syncInfos := syncm.SyncInfos{}
	syncInfos.Load()
	return syncm.BackupFiles(syncInfos, config.Conf.General.BackupFile)
}
