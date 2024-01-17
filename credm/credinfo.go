package credm

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
	CRED_TBL = "credentials"
)

type Credentials struct {
	Db          *db.Storage
	Logger      *zerolog.Logger
	Credentials []CredInfo `yaml:"credentials"`
}

type CredInfo struct {
	Application string `yaml:"application"`
	Username    string `yaml:"username"`
	Password    string `yaml:"-"`
	Timestamp   int64  `yaml:"timestamp"`
	Hash        string `yaml:"-"`
	// Not supported for now
	Tag []string `yaml:"tag"`
}

func (credentials *Credentials) Load() (map[interface{}]interface{}, error) {
	c, err := credentials.Db.GetFirst(CRED_TBL)
	if err != nil {
		c = map[interface{}]interface{}{}
	}

	credMap := c.(map[interface{}]interface{})
	for keyIface, valIFace := range credMap {
		credInfo := CredInfo{}
		mapstructure.Decode(valIFace, &credInfo)
		credInfo.Hash = keyIface.(string)
		credentials.Credentials = append(credentials.Credentials, credInfo)
	}

	credentials.SortByTimestamp()
	return credMap, nil
}

// ToString returns the contents of yaml file.
func (credentials *Credentials) ToString() (string, error) {
	var buffer bytes.Buffer
	err := yaml.NewEncoder(&buffer).Encode(credentials)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to yaml string: %v", err)
	}
	return buffer.String(), nil
}

// Order credentials regarding SortBy option defined in config yaml
// Prefix "-" reverses the order, default is "recency", "+<expressions>" is the same as "<expression>"
func (credentials *Credentials) Order() {
	sortBy := config.Conf.General.SortBy
	switch {
	case sortBy == "command" || sortBy == "+command":
		sort.Sort(ByTimestamp(credentials.Credentials))
	case sortBy == "-command":
		sort.Sort(sort.Reverse(ByTimestamp(credentials.Credentials)))

	case sortBy == "-recency":
		credentials.Reverse()
	}
}

func (credentials *Credentials) Reverse() {
	for i, j := 0, len(credentials.Credentials)-1; i < j; i, j = i+1, j-1 {
		credentials.Credentials[i], credentials.Credentials[j] = credentials.Credentials[j], credentials.Credentials[i]
	}
}

func (credentials *Credentials) SortByTimestamp() {
	sort.Sort(ByTimestamp(credentials.Credentials))
}

type ByTimestamp []CredInfo

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp > a[j].Timestamp }
