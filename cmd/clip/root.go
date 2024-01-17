package clip

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	version    = "dev"
)

var RootCmdClip = &cobra.Command{
	Use:           "clip",
	Short:         "Simple clipboard manager.",
	Long:          `clip - Simple clipboard manager.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := RootCmdClip.Execute(); err != nil {
		fmt.Printf("Error executing clip cmd, %v", err)
		os.Exit(-1)
	}
}
