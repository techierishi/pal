package snip

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/snipm"
	palSync "github.com/techierishi/pal/sync"
	"github.com/techierishi/pal/util"
)

var newCmd = &cobra.Command{
	Use:   "new COMMAND",
	Short: "Create a new snippet",
	Long:  `Create a new snippet (default: $HOME/.config/pal/snippet.yaml)`,
	RunE:  newFunc,
}

func init() {
	RootCmdSnip.AddCommand(newCmd)
}

func newFunc(cmd *cobra.Command, args []string) (err error) {
	logger := logr.GetLogInstance()

	snipDb, err := snipm.GetSnipDb(logger)
	if err != nil {
		return err
	}
	defer snipDb.Close()
	snippets := snipm.Snippets{
		Db:     snipDb,
		Logger: logger,
	}

	snippetMap, err := snippets.Load()
	if err != nil {
		return err
	}
	newSnippet := snipm.NewSnippet()
	newSnippet.Timestamp = util.UnixMilli()
	hash := util.CalculateHash(newSnippet.Command)
	if _, ok := snippetMap[hash]; ok {
		return fmt.Errorf("Snippet already exists")
	}
	snipDb.Upsert(fmt.Sprintf("%s.%s", snipm.SNIPPET_TBL, hash), newSnippet)

	fmt.Fprintf(color.Output, "%12s", color.GreenString("Snippet saved! \n"))

	if config.Conf.Gist.AutoSync {
		return palSync.AutoSync(config.Conf.General.BackupFile)
	}

	return nil
}
