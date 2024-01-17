package util

import "strings"

type CleanStr string

func (cstr CleanStr) StandardizeSpaces() CleanStr {
	return CleanStr(strings.Join(strings.Fields(string(cstr)), " "))
}

func (cstr CleanStr) ReplaceNewLine() CleanStr {
	lstr := strings.ReplaceAll(string(cstr), "\n", "#")
	return CleanStr(lstr)
}

func (cstr CleanStr) TruncateText(max int) CleanStr {
	var lstr string
	if len(string(cstr)) < max {
		lstr = string(cstr)
		return CleanStr(lstr)
	}

	lstr = string(cstr)[:max]
	return CleanStr(lstr)
}
