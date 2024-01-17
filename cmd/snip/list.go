package snip

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/snipm"
)

const (
	column = 40
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all snippets",
	Long:  `Show all snippets`,
	RunE:  listFunc,
}

func init() {
	RootCmdSnip.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&config.Flag.OneLine, "oneline", "", false,
		`Display snippets in one line`)
}

func listFunc(cmd *cobra.Command, args []string) error {
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

	if _, err := snippets.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, snippet := range snippets.Snippets {
		if config.Flag.OneLine {
			description := runewidth.FillRight(runewidth.Truncate(snippet.Description, col, "..."), col)
			command := runewidth.Truncate(snippet.Command, 100-4-col, "...")
			// make sure multiline command printed as oneline
			command = strings.Replace(command, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(description), color.YellowString(command))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("Description:"), snippet.Description)
			if strings.Contains(snippet.Command, "\n") {
				lines := strings.Split(snippet.Command, "\n")
				firstLine, restLines := lines[0], lines[1:]
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), firstLine)
				for _, line := range restLines {
					fmt.Fprintf(color.Output, "%12s %s\n",
						" ", line)
				}
			} else {
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.YellowString("    Command:"), snippet.Command)
			}
			if snippet.Tag != nil {
				tag := strings.Join(snippet.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("        Tag:"), tag)
			}
			if snippet.Output != "" {
				output := strings.Replace(snippet.Output, "\n", "\n             ", -1)
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.RedString("     Output:"), output)
			}
			fmt.Println(strings.Repeat("-", 30))
		}
	}
	return nil
}
