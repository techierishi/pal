package snip

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/syncm"
	"github.com/techierishi/pal/util"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit snippet file",
	Long:  `Edit snippet file (default: opened by vim)`,
	RunE:  editFunc,
}

func init() {
	RootCmdSnip.AddCommand(editCmd)
}

func editFunc(cmd *cobra.Command, args []string) (err error) {
	editor := config.Conf.General.Editor
	snippetFile := config.Conf.General.SnippetFile

	// file content before editing
	before := fileContent(snippetFile)

	err = util.EditFile(editor, snippetFile)
	if err != nil {
		return
	}

	// file content after editing
	after := fileContent(snippetFile)

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
