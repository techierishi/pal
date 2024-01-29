package alias

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmdAlias = &cobra.Command{
	Use:           "alias",
	Short:         "Simple alias manager.",
	Long:          `alias - Simple alias manager.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := RootCmdAlias.Execute(); err != nil {
		fmt.Printf("Error executing alias cmd, %v", err)
		os.Exit(-1)
	}
}
