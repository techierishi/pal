package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/techierishi/pal/aliasm"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/util"
)

var aliasInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initial pal",
	RunE:  initPalAppFunc,
}

func init() {
	RootCmd.AddCommand(aliasInitCmd)
}

func writeToBashrc() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}

	bashrcPath := homeDir + "/.bashrc"
	if _, err := os.Stat(bashrcPath); os.IsNotExist(err) {
		if err := os.WriteFile(bashrcPath, nil, 0644); err != nil {
			fmt.Println("Error creating .bashrc:", err)
			return
		}
	}

	if contains, err := util.FileContains(bashrcPath, util.PalShellrc); err != nil {
		fmt.Println("Error checking .bashrc:", err)
		return
	} else if !contains {
		if err := util.AppendToFile(bashrcPath, "\n"+util.PalShellrc+"\n"); err != nil {
			fmt.Println("Error adding pal shellrc to .bashrc:", err)
			return
		}
	}
}

func writeToZshrc() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return
	}

	zshrcPath := homeDir + "/.zshrc"

	if _, err := os.Stat(zshrcPath); os.IsNotExist(err) {
		if err := os.WriteFile(zshrcPath, nil, 0644); err != nil {
			fmt.Println("Error creating .zshrc:", err)
			return
		}
	}

	if contains, err := util.FileContains(zshrcPath, util.PalShellrc); err != nil {
		fmt.Println("Error checking .zshrc:", err)
		return
	} else if !contains {
		if err := util.AppendToFile(zshrcPath, "\n"+util.PalShellrc+"\n"); err != nil {
			fmt.Println("Error adding pal shellrc to .zshrc:", err)
			return
		}
	}

}

// Generate default Palrc file first
func initPalAppFunc(cmd *cobra.Command, args []string) (err error) {
	logger := logr.GetLogInstance()
	aliasDb, err := aliasm.GetAliasDb(logger)
	defer aliasDb.Close()
	aliases := aliasm.Aliases{
		Db: aliasDb,
	}
	if _, err = aliases.Load(); err != nil {
		return err
	}

	aliasStr, err := aliases.ToAliasString()
	if err != nil {
		return err
	}
	err = util.GeneratePalrc(aliasStr, true)
	if err != nil {
		return err
	}

	shell := util.GetCurrentShell()
	if strings.EqualFold(shell, "bash") {
		writeToBashrc()
	} else if strings.EqualFold(shell, "zsh") {
		writeToZshrc()
	}

	return nil
}
