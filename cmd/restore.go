package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	palSync "github.com/techierishi/pal/sync"
)

var restorePath string

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Backup all configs",
	Long:  `restore - Backup all configs to passed dir`,
	RunE:  restoreFunc,
}

func init() {
	RootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().StringVarP(&restorePath, "path", "p", "",
		`Path to restore from`)

}

func restoreFunc(cmd *cobra.Command, args []string) (err error) {
	dir, err := config.GetDefaultConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	if len(restorePath) <= 0 {
		fmt.Fprintf(os.Stderr, "Please pass path to restore file using `--path`")
		os.Exit(0)
	}
	return palSync.RestoreFiles(restorePath, dir)
}
