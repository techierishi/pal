package svc

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmdSvc = &cobra.Command{
	Use:           "svc",
	Short:         "Pal backgroud service.",
	Long:          `svc - Pal background service.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := RootCmdSvc.Execute(); err != nil {
		fmt.Printf("Error executing svc cmd, %v", err)
		os.Exit(-1)
	}
}
