package snip

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/snipm"
	"github.com/techierishi/pal/tui"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/alessio/shellescape.v1"
)

var delimiter string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search snippets",
	Long:  `Search snippets interactively (default filtering tool: <inbuilt>)`,
	RunE:  searchFunc,
}

func init() {
	RootCmdSnip.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	searchCmd.Flags().StringVarP(&config.Flag.Delimiter, "delimiter", "d", "; ",
		`Use delim as the command delimiter character`)
}

func searchFunc(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}

	customLabel := tui.CustomLabel{
		SearchTitle:   "Snippets",
		EnterHelpText: "print to terminal",
	}
	commands, err := snipm.FilterText(customLabel, options, flag.FilterTag)
	if err != nil {
		return err
	}

	fmt.Print(strings.Join(commands, flag.Delimiter))
	if terminal.IsTerminal(1) {
		fmt.Print("\n")
	}
	return nil
}
