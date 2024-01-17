package cred

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	version    = "dev"
)

var RootCmdCred = &cobra.Command{
	Use:           "cred",
	Short:         "Simple credential manager.",
	Long:          `cred - Simple credential manager.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := RootCmdCred.Execute(); err != nil {
		fmt.Printf("Error executing cred cmd, %v", err)
		os.Exit(-1)
	}
}
