package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/fatih/color"
	"github.com/techierishi/pal/logr"
)

// Generate only supports bash and zsh for now
var palrcContent string = `
####### DYNAMICALLY GENERATED FILE - DO NOT EDIT #######

####### <Alias> #######
{{{alias_lines}}}
alias "ph=pal hist list"
alias "pl=pal clip list"
alias "pe=pal snip exec --command"
alias "pr=pal cred list"
####### <Alias/> #######


####### <Bind> #######
{{{bind_lines}}}
####### <Bind/> #######
`

func GetBindList() (string, error) {
	var b strings.Builder
	shell := GetCurrentShell()

	if strings.EqualFold(shell, "cmd") || strings.EqualFold(shell, "bash") {
		b.WriteRune('\n')
		b.WriteString(`bind '"\C-r": "\C-a pal hist list -- \C-j"'`)
		b.WriteRune('\n')
	} else if strings.EqualFold(shell, "zsh") {
		b.WriteRune('\n')
		b.WriteString(`pal_ctrlr_widget() {
	pal hist list
}
zle -N pal_ctrlr_widget
bindkey "^R" pal_ctrlr_widget`)

		b.WriteRune('\n')
	}

	return b.String(), nil
}

func WriteStringToFile(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func GeneratePalrc(aliasLines string, showMessage bool) error {
	logger := logr.GetLogInstance()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Error().Any("Could not get user home directory while generating alias", err)
		return err
	}

	bindLines, err := GetBindList()
	data, err := mustache.Render(string(palrcContent),
		map[string]string{
			"alias_lines": aliasLines,
			"bind_lines":  bindLines,
		})

	if err != nil {
		logger.Error().Any("Could not replace mustache template", err)
		return err
	}

	err = WriteStringToFile(filepath.Join(homeDir, ".palrc"), data)
	if err != nil {
		logger.Fatal().AnErr("Could not write aliases", err)
		return err
	}

	if showMessage {
		fmt.Fprintf(color.Output, "%12s", color.GreenString(".palrc generated! \n"))
	}

	shell := GetCurrentShell()
	if strings.EqualFold(shell, "bash") {
		fmt.Fprintf(color.Output, "%12s", color.CyanString("Pleas run `source ~/.bashrc`\n"))
	} else if strings.EqualFold(shell, "zsh") {
		fmt.Fprintf(color.Output, "%12s", color.CyanString("Pleas run `source ~/.zshrc`\n"))
	}

	return nil
}
