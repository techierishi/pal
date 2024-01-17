package cred

import (
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/credm"
)

var credSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Password search",
	RunE:  searchFunc,
}

func init() {
	RootCmdCred.AddCommand(credSearchCmd)
}

func searchFunc(cmd *cobra.Command, args []string) (err error) {
	err = credm.PassSearch()
	return err
}
