package histm

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/tui"
	"github.com/techierishi/pal/util"
)

func histList() ([]list.Item, Histlist) {
	logger := logr.GetLogInstance()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error().Any("Could not get user home directory", err)
	}
	shell := util.GetCurrentShell()

	histStrList := Histlist{}
	if strings.EqualFold(shell, "cmd") || strings.EqualFold(shell, "bash") {
		bashHistoryPath := filepath.Join(homeDir, ".bash_history")
		histStrList = LoadCmdLinesFromBashFile(bashHistoryPath)
	} else if strings.EqualFold(shell, "zsh") {
		zshHistoryPath := filepath.Join(homeDir, ".zsh_history")
		histStrList = LoadCmdLinesFromZshFile(zshHistoryPath)
	}
	histStrList.Reverse()

	searchRowItems := make([]list.Item, 0)

	for idx, val := range histStrList.List {
		searchRowItems = append(searchRowItems, tui.NewSearchRowItem(val, strconv.Itoa(idx)))
	}

	return searchRowItems, histStrList

}

func HistList() (string, error) {
	searchRowItems, histStrList := histList()
	enterHelpText := "execute command"
	if config.Flag.Copy {
		enterHelpText = "copy to clipboard"

	}
	customLabel := tui.CustomLabel{
		SearchTitle:   "Shell history",
		EnterHelpText: enterHelpText,
	}
	selectedItem, err := tui.SearchUI(customLabel, searchRowItems)
	if err != nil {
		return "", err
	}
	idx, err := strconv.Atoi(selectedItem.Index())

	if len(histStrList.List) <= idx {
		return "", errors.New("bad index in history list")
	}
	return histStrList.List[idx], err
}
