package sync

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/syncm"
	"github.com/techierishi/pal/util"
)

// syncEditCmd represents the syncEdit command
var syncEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit sync config file",
	Long:  `Edit sync config (default: opened by vim)`,
	RunE:  syncEditFunc,
}

func init() {
	RootCmdSync.AddCommand(syncEditCmd)
}

func syncEditFunc(cmd *cobra.Command, args []string) (err error) {
	syncEditor := config.Conf.General.Editor
	syncFile := config.Conf.General.SyncFile

	// file content before syncEditing
	before := fileContent(syncFile)

	err = util.EditFile(syncEditor, syncFile)
	if err != nil {
		return
	}

	// file content after syncEditing
	after := fileContent(syncFile)

	// return if same file content
	if before == after {
		return nil
	}

	if config.Conf.Gist.AutoSync {
		return syncm.AutoSync(config.Conf.General.BackupFile)
	}

	return nil
}

func fileContent(fname string) string {
	data, _ := os.ReadFile(fname)
	return string(data)
}
