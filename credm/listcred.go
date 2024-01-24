package credm

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/tui"
	"github.com/techierishi/pal/util"
	"github.com/techierishi/pal/wrapper"
	"golang.design/x/clipboard"
)

func passList(credentials Credentials) []list.Item {
	passListItems := make([]list.Item, 0)

	idx := 0
	for _, credInfo := range credentials.Credentials {
		rowStr := fmt.Sprintf("[%s] %s", credInfo.Application, credInfo.Username)
		passListItems = append(passListItems, tui.NewSearchRowItem(rowStr, credInfo.Hash))
		idx++
	}

	return passListItems
}

func PassSearch() error {
	logger := logr.GetLogInstance()
	credDb, err := GetCredDb(logger)
	defer credDb.Close()
	if err != nil {
		return err
	}
	credentials := Credentials{
		Db:     credDb,
		Logger: logger,
	}
	credMap, err := credentials.Load()
	if err != nil {
		return err
	}
	passListItems := passList(credentials)
	customLabel := tui.CustomLabel{
		SearchTitle:   "Credentials",
		EnterHelpText: "copy to clipboard",
	}
	selectedItem, err := tui.SearchUI(customLabel, passListItems)
	if err != nil {
		return err
	}

	selectedCred := CredInfo{}
	mapstructure.Decode(credMap[selectedItem.Index()], &selectedCred)

	keyRing := wrapper.KeyRing{Logger: logger}
	passwd, ok := keyRing.Get(selectedCred.Application, selectedCred.Username)
	if !ok {
		fmt.Printf("%s\n", color.RedString("Keyring not supported on this os!"))
	}

	// Update the timestamp to show the mostly used item on top
	selectedCred.Timestamp = util.UnixMilli()
	credDb.Upsert(fmt.Sprintf("%s.%s", CRED_TBL, selectedItem.Index()), selectedCred)

	if config.Flag.HasClipboard {
		clipboard.Write(0, []byte(passwd))
		fmt.Printf("%s\n", color.GreenString("Copied credential to clipboard!"))
	} else {
		tui.PasswordModal(passwd)
	}

	return nil
}
