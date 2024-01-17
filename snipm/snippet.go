package snipm

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
	"github.com/ulfox/dby/db"
	"gopkg.in/yaml.v3"
)

var (
	SNIPPET_TBL = "snippets"
)

type Snippets struct {
	Db       *db.Storage
	Logger   *zerolog.Logger
	Snippets []SnippetInfo `yaml:"snippets"`
}

type SnippetInfo struct {
	Description string `yaml:"description"`
	Command     string `yaml:"command"`
	Timestamp   int64  `yaml:"timestamp"`
	Hash        string `yaml:"-"`
	// Not supported for now
	Tag        []string `yaml:"tag"`
	Output     string   `yaml:"output"`
	PromptPass string   `yaml:"promptPass"`
}

// Load reads yaml file.
func (snippets *Snippets) Load() (map[interface{}]interface{}, error) {
	c, err := snippets.Db.GetFirst(SNIPPET_TBL)
	if err != nil {
		c = map[interface{}]interface{}{}
	}
	snipMap := c.(map[interface{}]interface{})

	for keyIface, valIFace := range snipMap {
		snippetInfo := SnippetInfo{}
		mapstructure.Decode(valIFace, &snippetInfo)
		snippetInfo.Hash = keyIface.(string)
		snippets.Snippets = append(snippets.Snippets, snippetInfo)
	}

	snippets.SortByTimestamp()

	return snipMap, nil
}

// ToString returns the contents of yaml file.
func (snippets *Snippets) ToString() (string, error) {
	var buffer bytes.Buffer
	err := yaml.NewEncoder(&buffer).Encode(snippets)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to yaml string: %v", err)
	}
	return buffer.String(), nil
}

// Order snippets regarding SortBy option defined in config yaml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (snippets *Snippets) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "command" || sortBy == "+command":
		sort.Sort(ByCommand(snippets.Snippets))
	case sortBy == "-command":
		sort.Sort(sort.Reverse(ByCommand(snippets.Snippets)))

	case sortBy == "description" || sortBy == "+description":
		sort.Sort(ByDescription(snippets.Snippets))
	case sortBy == "-description":
		sort.Sort(sort.Reverse(ByDescription(snippets.Snippets)))

	case sortBy == "output" || sortBy == "+output":
		sort.Sort(ByOutput(snippets.Snippets))
	case sortBy == "-output":
		sort.Sort(sort.Reverse(ByOutput(snippets.Snippets)))

	case sortBy == "-recency":
		snippets.reverse()
	}
}

func (snippets *Snippets) reverse() {
	for i, j := 0, len(snippets.Snippets)-1; i < j; i, j = i+1, j-1 {
		snippets.Snippets[i], snippets.Snippets[j] = snippets.Snippets[j], snippets.Snippets[i]
	}
}

type ByCommand []SnippetInfo

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command > a[j].Command }

type ByDescription []SnippetInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Description > a[j].Description }

type ByOutput []SnippetInfo

func (a ByOutput) Len() int           { return len(a) }
func (a ByOutput) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOutput) Less(i, j int) bool { return a[i].Output > a[j].Output }

func (snippets *Snippets) SortByTimestamp() {
	sort.Sort(ByTimestamp(snippets.Snippets))
}

type ByTimestamp []SnippetInfo

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp > a[j].Timestamp }
