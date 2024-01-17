package clipm

import (
	"sort"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
)

var (
	CLIPBOARD_TBL = "clipboard"
)

type ClipInfos struct {
	Logger    *zerolog.Logger
	ClipInfos []ClipInfo `yaml:"clipinfos"`
}

type ClipInfo struct {
	Application string `yaml:"application"`
	Timestamp   int64  `yaml:"timestamp"`
	Content     string `yaml:"content"`
	Hash        string `yaml:"-"`
	// Not supported for now
	Tag []string `yaml:"tag"`
}

func (clipInfos *ClipInfos) Load(c interface{}) error {
	credMap := c.(map[interface{}]interface{})

	for keyIface, valIFace := range credMap {
		clipInfo := ClipInfo{}
		mapstructure.Decode(valIFace, &clipInfo)
		clipInfo.Hash = keyIface.(string)
		clipInfos.ClipInfos = append(clipInfos.ClipInfos, clipInfo)
	}

	clipInfos.SortByTimestamp()
	return nil
}

func (clipInfos *ClipInfos) reverse() {
	for i, j := 0, len(clipInfos.ClipInfos)-1; i < j; i, j = i+1, j-1 {
		clipInfos.ClipInfos[i], clipInfos.ClipInfos[j] = clipInfos.ClipInfos[j], clipInfos.ClipInfos[i]
	}
}

func (clipInfos *ClipInfos) SortByTimestamp() {
	sort.Sort(ByTimestamp(clipInfos.ClipInfos))
}

type ByTimestamp []ClipInfo

func (a ByTimestamp) Len() int           { return len(a) }
func (a ByTimestamp) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimestamp) Less(i, j int) bool { return a[i].Timestamp > a[j].Timestamp }
