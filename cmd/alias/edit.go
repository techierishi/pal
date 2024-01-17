package alias

import (
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/aliasm"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/util"
)

var credEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit alias file",
	Long:  `Edit alias file (default: opened by vim)`,
	RunE:  editFunc,
}

func init() {
	RootCmdAlias.AddCommand(credEditCmd)
}

func editFunc(cmd *cobra.Command, args []string) (err error) {
	logger := logr.GetLogInstance()

	editor := config.Conf.General.Editor
	aliasFile := config.Conf.General.AliasFile
	err = util.EditFile(editor, aliasFile)
	if err != nil {
		return
	}

	aliasDb, err := aliasm.GetAliasDb(logger)
	defer aliasDb.Close()
	aliases := aliasm.Aliases{
		Db:     aliasDb,
		Logger: logger,
	}
	if _, err = aliases.Load(); err != nil {
		return err
	}

	aliasStr, err := aliases.ToAliasString()
	if err != nil {
		return err
	}
	err = util.GeneratePalrc(aliasStr, true)
	if err != nil {
		return err
	}

	return nil
}
