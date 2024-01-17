package cred

import (
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/util"
)

var credEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit credential file",
	Long:  `Edit credential file (default: opened by vim)`,
	RunE:  editFunc,
}

func init() {
	RootCmdCred.AddCommand(credEditCmd)
}

func editFunc(cmd *cobra.Command, args []string) (err error) {
	editor := config.Conf.General.Editor
	credFile := config.Conf.General.CredFile
	err = util.EditFile(editor, credFile)
	if err != nil {
		return
	}

	return nil
}
