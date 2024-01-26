package alias

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/aliasm"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/syncm"
	"github.com/techierishi/pal/util"
)

var aliasNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create new alias",
	RunE:  newFunc,
}

func init() {
	RootCmdAlias.AddCommand(aliasNewCmd)
}

func newFunc(cmd *cobra.Command, args []string) (err error) {
	logger := logr.GetLogInstance()

	aliasDb, err := aliasm.GetAliasDb(logger)
	defer aliasDb.Close()
	aliases := aliasm.Aliases{
		Db:     aliasDb,
		Logger: logger,
	}
	aliasMap, err := aliases.Load()
	if err != nil {
		return err
	}

	newAlias := aliasm.NewAlias()
	newAlias.Timestamp = util.UnixMilli()
	hash := util.CalculateHash(newAlias.Command)
	if _, ok := aliasMap[hash]; ok {
		return fmt.Errorf("Alias already exists")
	}

	aliasDb.Upsert(fmt.Sprintf("%s.%s", aliasm.ALIAS_TBL, hash), newAlias)
	fmt.Fprintf(color.Output, "%12s", color.GreenString("Alias saved! \n"))

	aliasStr, err := aliases.ToAliasString()
	if err != nil {
		return err
	}
	err = util.GeneratePalrc(aliasStr, true)
	if err != nil {
		return err
	}

	if config.Conf.Gist.AutoSync {
		return syncm.AutoSync(config.Conf.General.BackupFile)
	}

	return nil
}
