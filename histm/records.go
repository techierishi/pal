package histm

import (
	"bufio"
	"os"
	"strings"

	"github.com/techierishi/pal/logr"
)

// LoadCmdLinesFromZshFile loads cmdlines from zsh history file
func LoadCmdLinesFromZshFile(fname string) Histlist {
	logger := logr.GetLogInstance()

	hl := New()

	file, err := os.Open(fname)
	if err != nil {
		logger.Error().Any("Failed to open zsh history file - skipping reading zsh history", err)
		return hl
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// trim newline
		line = strings.TrimRight(line, "\n")
		var cmd string
		// zsh format EXTENDED_HISTORY
		// : 1576270617:0;make install
		// zsh format no EXTENDED_HISTORY
		// make install
		if len(line) == 0 {
			// skip empty
			continue
		}
		if strings.Contains(line, ":") && strings.Contains(line, ";") &&
			len(strings.Split(line, ":")) >= 3 && len(strings.Split(line, ";")) >= 2 {
			// contains at least 2x ':' and 1x ';' => assume EXTENDED_HISTORY
			cmd = strings.Split(line, ";")[1]
		} else {
			cmd = line
		}
		hl.AddCmdLine(cmd)
	}
	return hl
}

// LoadCmdLinesFromBashFile loads cmdlines from bash history file
func LoadCmdLinesFromBashFile(fname string) Histlist {
	logger := logr.GetLogInstance()

	hl := New()
	file, err := os.Open(fname)
	if err != nil {
		logger.Error().Any("Failed to open bash history file - skipping reading bash history", err)
		return hl
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// trim newline
		line = strings.TrimRight(line, "\n")
		// trim spaces from left
		line = strings.TrimLeft(line, " ")
		// bash format (two lines)
		// #1576199174
		// make install
		if strings.HasPrefix(line, "#") {
			// is either timestamp or comment => skip
			continue
		}
		if len(line) == 0 {
			// skip empty
			continue
		}
		hl.AddCmdLine(line)
	}
	return hl
}
