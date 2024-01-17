package aliasm

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
	"github.com/ulfox/dby/db"
	"gopkg.in/yaml.v3"
)

var (
	ALIAS_TBL = "aliases"
)

type Aliases struct {
	Db      *db.Storage
	Logger  *zerolog.Logger
	Aliases []AliasInfo `yaml:"aliases"`
}

type AliasInfo struct {
	Alias     string `yaml:"alias"`
	Command   string `yaml:"command"`
	Timestamp int64  `yaml:"timestamp"`
	Hash      string `yaml:"-"`
	// Not supported for now
	Tag []string `yaml:"tag"`
}

func (alias *Aliases) Load() (map[interface{}]interface{}, error) {
	c, err := alias.Db.GetFirst(ALIAS_TBL)
	if err != nil {
		c = map[interface{}]interface{}{}
	}
	aliasMap := c.(map[interface{}]interface{})

	for keyIface, valIFace := range aliasMap {
		aliasInfo := AliasInfo{}
		mapstructure.Decode(valIFace, &aliasInfo)
		aliasInfo.Hash = keyIface.(string)
		alias.Aliases = append(alias.Aliases, aliasInfo)
	}
	return aliasMap, nil
}

// ToString returns the contents of yaml file.
func (aliases *Aliases) ToString() (string, error) {
	var buffer bytes.Buffer
	err := yaml.NewEncoder(&buffer).Encode(aliases)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to yaml string: %v", err)
	}
	return buffer.String(), nil
}

func (alias *Aliases) ToAliasString() (string, error) {
	var b strings.Builder
	for _, alias := range alias.Aliases {

		if strings.Contains(alias.Command, "\n") {
			lines := strings.Split(alias.Command, "\n")

			b.WriteRune('\n')
			b.WriteString(fmt.Sprintf("function %s {", alias.Alias))
			b.WriteRune('\n')
			for _, line := range lines {
				b.WriteString(fmt.Sprintf(" %s", line))
				b.WriteRune('\n')
			}
			b.WriteString("}")
			b.WriteRune('\n')
		} else {
			b.WriteRune('\n')
			b.WriteString(fmt.Sprintf("alias \"%s=%s\"", alias.Alias, alias.Command))
			b.WriteRune('\n')
		}
	}

	return b.String(), nil
}

// Order aliases regarding SortBy option defined in config yaml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (aliases *Aliases) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "command" || sortBy == "+command":
		sort.Sort(ByCommand(aliases.Aliases))
	case sortBy == "-command":
		sort.Sort(sort.Reverse(ByCommand(aliases.Aliases)))

	case sortBy == "alias" || sortBy == "+alias":
		sort.Sort(ByDescription(aliases.Aliases))
	case sortBy == "-alias":
		sort.Sort(sort.Reverse(ByDescription(aliases.Aliases)))

	case sortBy == "-recency":
		aliases.reverse()
	}
}

func (aliases *Aliases) reverse() {
	for i, j := 0, len(aliases.Aliases)-1; i < j; i, j = i+1, j-1 {
		aliases.Aliases[i], aliases.Aliases[j] = aliases.Aliases[j], aliases.Aliases[i]
	}
}

type ByCommand []AliasInfo

func (a ByCommand) Len() int           { return len(a) }
func (a ByCommand) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCommand) Less(i, j int) bool { return a[i].Command > a[j].Command }

type ByDescription []AliasInfo

func (a ByDescription) Len() int           { return len(a) }
func (a ByDescription) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDescription) Less(i, j int) bool { return a[i].Alias > a[j].Alias }
