package snip

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/snipm"
	"github.com/techierishi/pal/tui"
	"golang.design/x/clipboard"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy the selected commands",
	Long:  `Copy the selected commands to clipboard`,
	RunE:  copyFunc,
}

func init() {
	RootCmdSnip.AddCommand(copyCmd)
	// Query not working currently
	copyCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	copyCmd.Flags().BoolVarP(&config.Flag.Command, "command", "", false,
		`Display snippets in one line`)
	copyCmd.Flags().StringVarP(&config.Flag.Delimiter, "delimiter", "d", "; ",
		`Use delim as the command delimiter character`)
}

func copyFunc(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", flag.Query))
	}
	customLabel := tui.CustomLabel{
		SearchTitle:   "Snippets",
		EnterHelpText: "copy to clipboard",
	}
	commands, err := snipm.FilterText(customLabel, options, flag.FilterTag)
	if err != nil {
		return err
	}
	command := strings.Join(commands, flag.Delimiter)
	if flag.Command && command != "" {
		fmt.Printf("%s: %s\n", color.YellowString("Command"), command)
	}

	if config.Flag.HasClipboard {
		clipboard.Write(0, []byte(command))
		fmt.Printf("%s\n", color.GreenString("Copied command to clipboard!"))
	} else {
		fmt.Printf("%s\n", color.RedString("Clipboard API not available in this system!"))
	}

	return nil
}
