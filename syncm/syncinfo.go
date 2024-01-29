package syncm

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"gopkg.in/yaml.v3"
)

type SyncInfos struct {
	Files []SyncInfo `yaml:"files"`
}

type SyncInfo struct {
	FilePath string            `yaml:"filepath,omitempty"`
	Metadata map[string]string `yaml:"metadata"`
}

// Load reads yaml file.
func (syncInfos *SyncInfos) Load() error {
	logger := logr.GetLogInstance()

	// Adding pal files by default
	syncInfos.Files = append(syncInfos.Files, SyncInfo{
		FilePath: config.Conf.General.SnippetFile,
	})
	syncInfos.Files = append(syncInfos.Files, SyncInfo{
		FilePath: config.Conf.General.CredFile,
	})
	syncInfos.Files = append(syncInfos.Files, SyncInfo{
		FilePath: config.Conf.General.AliasFile,
	})
	syncInfos.Files = append(syncInfos.Files, SyncInfo{
		FilePath: config.Conf.General.SyncFile,
	})

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Load:: Error getting user home directory:", err)
		return err
	}

	confDir, err := config.GetUserConfigDir()
	if err != nil {
		fmt.Println("Load:: Error getting user config directory:", err)
		return err
	}

	syncInfoFilePath := config.Conf.General.SyncFile
	if _, err := os.Stat(syncInfoFilePath); os.IsNotExist(err) {
		return nil
	}
	syncInfoFile, err := os.OpenFile(syncInfoFilePath, os.O_RDONLY, 0600)
	if err != nil {
		logger.Error().Any("error opening/creating file: %v", err)
		return err
	}
	defer syncInfoFile.Close()
	dec := yaml.NewDecoder(syncInfoFile)
	err = dec.Decode(&syncInfos)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		logger.Error().Any("Failed to load syncInfo file. %v", err)
		return err
	}

	// replace signs
	for idx, fileInfo := range syncInfos.Files {
		fullPath, dirCheck := ReplaceDirFromSign(fileInfo.FilePath, homeDir, confDir)
		if !dirCheck {
			continue
		}
		syncInfos.Files[idx].FilePath = fullPath
	}

	syncInfos.Files = removeDuplicate(syncInfos.Files)

	syncInfos.Order()
	return nil
}

// Save saves the syncInfos to yaml file.
func (syncInfos *SyncInfos) Save() error {
	syncInfo := config.Conf.General.BackupFile
	f, err := os.Create(syncInfo)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save syncInfo file. err: %s", err)
	}
	return yaml.NewEncoder(f).Encode(syncInfos)
}

func (syncInfos *SyncInfos) ToString() (string, error) {
	var buffer bytes.Buffer
	err := yaml.NewEncoder(&buffer).Encode(syncInfos)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to yaml string: %v", err)
	}
	return buffer.String(), nil
}

func (syncInfos *SyncInfos) Order() {
	sortBy := config.Conf.General.SortBy
	switch {

	case sortBy == "-recency":
		syncInfos.Reverse()
	}
}

func (syncInfos *SyncInfos) Reverse() {
	for i, j := 0, len(syncInfos.Files)-1; i < j; i, j = i+1, j-1 {
		syncInfos.Files[i], syncInfos.Files[j] = syncInfos.Files[j], syncInfos.Files[i]
	}
}

func removeDuplicate(sliceList []SyncInfo) []SyncInfo {
	allKeys := make(map[string]bool)
	list := []SyncInfo{}
	for _, item := range sliceList {
		if _, value := allKeys[item.FilePath]; !value {
			allKeys[item.FilePath] = true
			list = append(list, item)
		}
	}
	return list
}
