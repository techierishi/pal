package cred

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/credm"
	"github.com/techierishi/pal/logr"
)

const (
	column = 40
)

var credListCmd = &cobra.Command{
	Use:   "list",
	Short: "Password list to console",
	RunE:  listFunc,
}

func init() {
	RootCmdCred.AddCommand(credListCmd)
}

func listFunc(cmd *cobra.Command, args []string) (err error) {
	logger := logr.GetLogInstance()
	credDb, err := credm.GetCredDb(logger)
	defer credDb.Close()
	if err != nil {
		return err
	}
	credentials := credm.Credentials{
		Db:     credDb,
		Logger: logger,
	}
	if _, err := credentials.Load(); err != nil {
		return err
	}

	col := config.Conf.General.Column
	if col == 0 {
		col = column
	}

	for _, credential := range credentials.Credentials {
		if config.Flag.OneLine {
			application := runewidth.FillRight(runewidth.Truncate(credential.Application, col, "..."), col)
			username := runewidth.Truncate(credential.Username, 100-4-col, "...")
			username = strings.Replace(username, "\n", "\\n", -1)
			fmt.Fprintf(color.Output, "%s : %s\n",
				color.GreenString(application), color.YellowString(username))
		} else {
			fmt.Fprintf(color.Output, "%12s %s\n",
				color.GreenString("Application:"), credential.Application)

			fmt.Fprintf(color.Output, "%12s %s\n",
				color.YellowString("   Username:"), credential.Username)

			if credential.Tag != nil {
				tag := strings.Join(credential.Tag, " ")
				fmt.Fprintf(color.Output, "%12s %s\n",
					color.CyanString("        Tag:"), tag)
			}
			fmt.Println(strings.Repeat("-", 30))
		}
	}
	return nil
}
