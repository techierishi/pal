package sync

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/syncm"
)

var RootCmdSync = &cobra.Command{
	Use:           "sync",
	Short:         "Sync configs",
	Long:          `Sync configs with gist/gitlab`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          syncFunc,
}

func Execute() {
	if err := RootCmdSync.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func syncFunc(cmd *cobra.Command, args []string) (err error) {
	syncInfos := syncm.SyncInfos{}
	syncInfos.Load()
	err = syncm.BackupFiles(syncInfos, config.Conf.General.BackupFile)
	if err != nil {
		return
	}

	return syncm.AutoSync(config.Conf.General.BackupFile)
}
