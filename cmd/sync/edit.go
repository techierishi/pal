package sync

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/syncm"
	"github.com/techierishi/pal/util"
)

var syncFileExample bool

// syncEditCmd represents the syncEdit command
var syncEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit sync config file",
	Long:  `Edit sync config (default: opened by vim)`,
	RunE:  syncEditFunc,
}

func init() {
	RootCmdSync.AddCommand(syncEditCmd)
	syncEditCmd.PersistentFlags().BoolVarP(&syncFileExample, "example", "", false, "Show sync file example")

}

func syncEditFunc(cmd *cobra.Command, args []string) (err error) {

	if syncFileExample {
		fmt.Fprintf(color.Output, "%s\n",
			color.CyanString("# Sync file only supports files from user home directory"))
		fmt.Fprintf(color.Output, "%s\n",
			color.CyanString("# Use <home> to denote home directory. Example follows"))
		fmt.Fprintf(color.Output, "%s\n",
			color.GreenString("files:"))
		fmt.Fprintf(color.Output, "%s\n",
			color.GreenString("  - filepath: <home>/.bashrc"))
		fmt.Fprintf(color.Output, "%s\n",
			color.GreenString("  - filepath: <home>/.config/nvim/init.vim"))
		return
	}

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
