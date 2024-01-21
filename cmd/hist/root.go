package hist

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	histm "github.com/techierishi/pal/histm"
	"github.com/techierishi/pal/util"
	"golang.design/x/clipboard"
)

var (
	configFile string
	version    = "dev"
)

var RootCmdHist = &cobra.Command{
	Use:           "hist",
	Short:         "Simple shell history manager.",
	Long:          `hist - Simple shell history manager.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          execFunc,
}

func init() {
	RootCmdHist.Flags().BoolVarP(&config.Flag.Command, "command", "c", false,
		`Show the command with the plain text before executing`)
	RootCmdHist.Flags().BoolVarP(&config.Flag.Copy, "copy", "p", false,
		`Just copy command with the plain text.`)
}

func execFunc(cmd *cobra.Command, args []string) (err error) {

	command, err := histm.HistList()
	if err != nil {
		return err
	}
	if config.Flag.Debug {
		fmt.Printf("Command: %s\n", command)
	}
	if config.Flag.Command {
		fmt.Printf("%s: %s\n", color.YellowString("Command"), command)
	}
	if config.Flag.Copy {
		if config.Flag.HasClipboard {
			clipboard.Write(0, []byte(command))
			fmt.Printf("%s\n", color.GreenString("Copied selected command!"))
		}
		return nil
	}

	return util.RunCmd(command, os.Stdin, os.Stdout)
}

func Execute() {
	if err := RootCmdHist.Execute(); err != nil {
		fmt.Printf("Error executing hist cmd, %v", err)
		os.Exit(-1)
	}
}
