package alias

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/aliasm"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
)

const (
	column = 40
)

var credListCmd = &cobra.Command{
	Use:   "list",
	Short: "Password list to console",
	RunE:  listFunc,
}

func init() {
	RootCmdAlias.AddCommand(credListCmd)
}

func listFunc(cmd *cobra.Command, args []string) (err error) {
	logger := logr.GetLogInstance()

	aliasDb, err := aliasm.GetAliasDb(logger)
	defer aliasDb.Close()
	aliases := aliasm.Aliases{
		Db:     aliasDb,
		Logger: logger,
	}
	if _, err = aliases.Load(); err != nil {
		return err
	}
	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, alias := range aliases.Aliases {
		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(alias.Alias, col, "..."), col)
			command := runewidth.Truncate(alias.Command, 100-4-col, "...")
			// make sure multiline command printed as oneline
			command = strings.Replace(command, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(description), color.YellowString(command))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("  Alias:"), alias.Alias)
			if strings.Contains(alias.Command, "\n") {
				lines := strings.Split(alias.Command, "\n")
				firstLine, restLines := lines[0], lines[1:]
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("Command:"), firstLine)
				for _, line := range restLines {
					fmt.Fprintf(color.Output, "%8s %s\n",
						"", line)
				}
			} else {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("Command:"), alias.Command)
			}
			if alias.Tag != nil {
				tag := strings.Join(alias.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("    Tag:"), tag)
			}
			fmt.Println(strings.Repeat("-", 30))
		}
	}
	return nil
}
