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
func (backups *SyncInfos) Load() error {
	logger := logr.GetLogInstance()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("RestoreFiles:: Error getting user home directory:", err)
		return err
	}

	confDir, err := config.GetUserConfigDir()
	if err != nil {
		fmt.Println("RestoreFiles:: Error getting user config directory:", err)
		return err
	}

	syncFilePath := config.Conf.General.SyncFile
	if _, err := os.Stat(syncFilePath); os.IsNotExist(err) {
		return nil
	}
	syncFile, err := os.OpenFile(syncFilePath, os.O_RDONLY, 0600)
	if err != nil {
		logger.Error().Any("error opening/creating file: %v", err)
		return err
	}
	defer syncFile.Close()
	dec := yaml.NewDecoder(syncFile)
	err = dec.Decode(&backups)
	if err != nil {
		if err == io.EOF {
			fmt.Println("SyncInfos file is empty.")
			return nil
		}
		return fmt.Errorf("Failed to load backup file. %v", err)
	}

	// replace signs
	for idx, fileInfo := range backups.Files {
		dirPath, dirCheck := ReplaceDirFromSign(fileInfo.FilePath, homeDir, confDir)
		if !dirCheck {
			continue
		}
		backups.Files[idx].FilePath = dirPath

	}

	backups.Files = append(backups.Files, SyncInfo{
		FilePath: config.Conf.General.SnippetFile,
	})
	backups.Files = append(backups.Files, SyncInfo{
		FilePath: config.Conf.General.CredFile,
	})
	backups.Files = append(backups.Files, SyncInfo{
		FilePath: config.Conf.General.AliasFile,
	})

	backups.Order()
	return nil
}

// Save saves the backups to yaml file.
func (backups *SyncInfos) Save() error {
	syncFile := config.Conf.General.BackupFile
	f, err := os.Create(syncFile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save backup file. err: %s", err)
	}
	return yaml.NewEncoder(f).Encode(backups)
}

func (backups *SyncInfos) ToString() (string, error) {
	var buffer bytes.Buffer
	err := yaml.NewEncoder(&buffer).Encode(backups)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to yaml string: %v", err)
	}
	return buffer.String(), nil
}

func (backups *SyncInfos) Order() {
	sortBy := config.Conf.General.SortBy
	switch {

	case sortBy == "-recency":
		backups.Reverse()
	}
}

func (backups *SyncInfos) Reverse() {
	for i, j := 0, len(backups.Files)-1; i < j; i, j = i+1, j-1 {
		backups.Files[i], backups.Files[j] = backups.Files[j], backups.Files[i]
	}
}
