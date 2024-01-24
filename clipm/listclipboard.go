package clipm

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/mitchellh/mapstructure"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/svcm"
	"github.com/techierishi/pal/tui"
	"github.com/techierishi/pal/util"
)

func clipItemList(clipData map[interface{}]interface{}) ([]list.Item, error) {

	clipListItems := make([]list.Item, 0)
	clipInfos := ClipInfos{}
	if err := clipInfos.Load(clipData); err != nil {
		return nil, err
	}

	for _, clipInfo := range clipInfos.ClipInfos {
		str := util.CleanStr(clipInfo.Content).StandardizeSpaces().TruncateText(150).ReplaceNewLine()
		clipListItems = append(clipListItems, tui.NewSearchRowItem(string(str), clipInfo.Hash))

	}

	return clipListItems, nil
}

func ClipboardList() (*ClipInfo, error) {

	logger := logr.GetLogInstance()
	svcDb, err := svcm.GetSvcDb(logger)
	if err != nil {
		logger.Error().Any("error", err).Msg("Error opening app db")
	}
	defer svcDb.Close()

	_, err = svcDb.GetPath("pid")
	if err != nil {
		fmt.Println("Please start pal server using `pal svc start -d` to record clipboard events.")
		os.Exit(0)
	}

	clipDb, err := GetClipDb(logger)
	if err != nil {
		return nil, err
	}
	defer clipDb.Close()

	c, err := clipDb.GetPath("clipboard")
	clipMap := map[interface{}]interface{}{}
	if err == nil {
		clipMap = c.(map[interface{}]interface{})
	}

	clipListItems, err := clipItemList(clipMap)
	if err != nil {
		return nil, err
	}

	customLabel := tui.CustomLabel{
		SearchTitle:   "Clipboard",
		EnterHelpText: "copy to clipboard",
	}
	selectedItem, err := tui.SearchUI(customLabel, clipListItems)

	if err != nil {
		return nil, err
	}

	clipInfo := ClipInfo{}
	selectedClipItem := clipMap[interface{}(selectedItem.Index())]
	mapstructure.Decode(selectedClipItem, &clipInfo)

	return &clipInfo, nil
}
