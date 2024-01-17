package snipm

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/credm"
	"github.com/techierishi/pal/dialog"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/search"
	"github.com/techierishi/pal/util"
)

func FilterText(customLabel search.CustomLabel, options []string, tag string) (commands []string, err error) {
	logger := logr.GetLogInstance()

	snipDb, err := GetSnipDb(logger)
	if err != nil {
		return nil, err
	}
	defer snipDb.Close()
	snippets := Snippets{
		Db:     snipDb,
		Logger: logger,
	}

	snippetMap, err := snippets.Load()
	if err != nil {
		return nil, err
	}

	if 0 < len(tag) {
		var filteredSnippets Snippets
		for _, snippet := range snippets.Snippets {
			for _, t := range snippet.Tag {
				if tag == t {
					filteredSnippets.Snippets = append(filteredSnippets.Snippets, snippet)
				}
			}
		}
		snippets = filteredSnippets
	}

	snippetTexts := map[string]SnippetInfo{}
	snipListItems := make([]list.Item, 0)

	for _, s := range snippets.Snippets {
		command := s.Command
		if strings.ContainsAny(command, "\n") {
			command = strings.Replace(command, "\n", "\\n", -1)
		}
		t := fmt.Sprintf("[%s]: %s", s.Description, command)

		tags := ""
		for _, tag := range s.Tag {
			tags += fmt.Sprintf(" #%s", tag)
		}
		t += tags

		snippetTexts[t] = s
		if config.Flag.Color {
			t = fmt.Sprintf("[%s]: %s%s",
				color.RedString(s.Description), command, color.BlueString(tags))
		}

		rowstr := util.CleanStr(t).StandardizeSpaces().TruncateText(300).ReplaceNewLine()
		snipListItems = append(snipListItems, search.NewSearchRowItem(string(rowstr), s.Hash))

	}

	selectedItem, err := search.SearchUI(customLabel, snipListItems)
	if err != nil {
		return nil, err
	}

	selectedSnippet := SnippetInfo{}
	mapstructure.Decode(snippetMap[selectedItem.Index()], &selectedSnippet)
	if err != nil {
		return commands, err
	}
	// Update the timestamp to show the mostly used item on top
	selectedSnippet.Timestamp = util.UnixMilli()
	snipDb.Upsert(fmt.Sprintf("%s.%s", SNIPPET_TBL, selectedItem.Index()), selectedSnippet)

	lines := strings.Split(strings.TrimSuffix(selectedSnippet.Command, "\n"), "\n")

	params := dialog.SearchForParams(lines)
	promptPass, err := util.ParseBool(selectedSnippet.PromptPass)
	if err != nil {
		promptPass = false
	}

	if params != nil {
		dialog.CurrentCommand = selectedSnippet.Command
		dialog.GenerateParamsForm(params, dialog.CurrentCommand)

		if promptPass {
			credm.PassSearch()
		}
		res := []string{dialog.FinalCommand}
		return res, nil
	}

	if promptPass {
		credm.PassSearch()
	}
	for _, line := range lines {
		commands = append(commands, fmt.Sprint(line))
	}

	return commands, nil
}
