package cmd

import (
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/aliasm"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/util"
)

var aliasGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate aliases, key mappings",
	RunE:  genFunc,
}

func init() {
	RootCmd.AddCommand(aliasGenCmd)
}

func genFunc(cmd *cobra.Command, args []string) (err error) {
	logger := logr.GetLogInstance()

	aliasDb, err := aliasm.GetAliasDb(logger)
	defer aliasDb.Close()
	aliases := aliasm.Aliases{
		Db: aliasDb,
	}
	if _, err = aliases.Load(); err != nil {
		return err
	}

	aliasStr, err := aliases.ToAliasString()
	if err != nil {
		return err
	}
	return util.GeneratePalrc(aliasStr, true)
}
