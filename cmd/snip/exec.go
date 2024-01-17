package snip

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/search"
	"github.com/techierishi/pal/snipm"
	"github.com/techierishi/pal/util"
	"gopkg.in/alessio/shellescape.v1"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Run the selected commands",
	Long:  `Run the selected commands directly`,
	RunE:  execFunc,
}

func init() {
	RootCmdSnip.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&config.Flag.Query, "query", "q", "",
		`Initial value for query`)
	execCmd.Flags().BoolVarP(&config.Flag.Command, "command", "c", false,
		`Show the command with the plain text before executing`)
}

func execFunc(cmd *cobra.Command, args []string) (err error) {
	flag := config.Flag

	var options []string
	if flag.Query != "" {
		options = append(options, fmt.Sprintf("--query %s", shellescape.Quote(flag.Query)))
	}
	customLabel := search.CustomLabel{
		SearchTitle:   "Snippets",
		EnterHelpText: "enter to execute",
	}
	commands, err := snipm.FilterText(customLabel, options, flag.FilterTag)
	if err != nil {
		return err
	}
	command := strings.Join(commands, "; ")
	if config.Flag.Debug {
		fmt.Printf("Command: %s\n", command)
	}
	if config.Flag.Command {
		fmt.Printf("%s: %s\n", color.YellowString("Command"), command)
	}
	return util.RunCmd(command, os.Stdin, os.Stdout)
}
