package snip

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	version    = "dev"
)

var RootCmdSnip = &cobra.Command{
	Use:           "snip",
	Short:         "Simple command-line snippet manager.",
	Long:          `snip - Simple command-line snippet manager.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := RootCmdSnip.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
