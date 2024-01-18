package credm

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/search"
	"github.com/techierishi/pal/util"
	"github.com/zalando/go-keyring"
	"golang.design/x/clipboard"
)

func passList(credentials Credentials) []list.Item {
	passListItems := make([]list.Item, 0)

	idx := 0
	for _, credInfo := range credentials.Credentials {
		rowStr := fmt.Sprintf("[%s] %s", credInfo.Application, credInfo.Username)
		passListItems = append(passListItems, search.NewSearchRowItem(rowStr, credInfo.Hash))
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
	customLabel := search.CustomLabel{
		SearchTitle:   "Credentials",
		EnterHelpText: "copy to clipboard",
	}
	selectedItem, err := search.SearchUI(customLabel, passListItems)
	if err != nil {
		return err
	}

	selectedCred := CredInfo{}
	mapstructure.Decode(credMap[selectedItem.Index()], &selectedCred)
	passwd, err := keyring.Get(selectedCred.Application, selectedCred.Username)
	if err != nil {
		logger.Error().AnErr("msg", err)
		return err
	}

	clipboard.Write(0, []byte(passwd))

	// Update the timestamp to show the mostly used item on top
	selectedCred.Timestamp = util.UnixMilli()
	credDb.Upsert(fmt.Sprintf("%s.%s", CRED_TBL, selectedItem.Index()), selectedCred)

	fmt.Printf("%s\n", color.GreenString("Copied credential to clipboard!"))

	return nil
}
