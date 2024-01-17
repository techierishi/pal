package cmd

import (
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/util"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Edit config file",
	Long:  `Edit config file (default: opened by vim)`,
	RunE:  configureFunc,
}

func init() {
	RootCmd.AddCommand(configureCmd)
}

func configureFunc(cmd *cobra.Command, args []string) (err error) {
	editor := config.Conf.General.Editor
	return util.EditFile(editor, configFile)
}
