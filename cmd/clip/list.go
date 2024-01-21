package clip

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/techierishi/pal/clipm"
	"github.com/techierishi/pal/config"
	clpbaord "golang.design/x/clipboard"
)

var clipListCmd = &cobra.Command{
	Use:   "list",
	Short: "Clipboard history",
	RunE:  listFunc,
}

func init() {
	RootCmdClip.AddCommand(clipListCmd)
}

func listFunc(cmd *cobra.Command, args []string) (err error) {
	selectedClipItem, err := clipm.ClipboardList()
	if err != nil {
		return err
	}

	if config.Flag.HasClipboard {
		clpbaord.Write(0, []byte(selectedClipItem.Content))
		fmt.Printf("%s\n", color.GreenString("Copied selected item!"))
	}

	return nil
}
