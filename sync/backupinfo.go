package sync

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"gopkg.in/yaml.v3"
)

type Backups struct {
	Files []FileInfo `yaml:"files"`
}

type FileInfo struct {
	FileContent string            `yaml:"filecontent,omitempty"`
	Metadata    map[string]string `yaml:"metadata"`
}

// Load reads yaml file.
func (backups *Backups) Load() error {
	logger := logr.GetLogInstance()

	bkpFilePath := config.Conf.General.BackupFile
	if _, err := os.Stat(bkpFilePath); os.IsNotExist(err) {
		return nil
	}
	backupFile, err := os.OpenFile(bkpFilePath, os.O_RDONLY, 0600)
	if err != nil {
		logger.Error().Any("error opening/creating file: %v", err)
		return err
	}
	defer backupFile.Close()
	dec := yaml.NewDecoder(backupFile)
	err = dec.Decode(&backups)
	if err != nil {
		if err == io.EOF {
			fmt.Println("Backups file is empty.")
			return nil
		}
		return fmt.Errorf("Failed to load backup file. %v", err)
	}
	backups.Order()
	return nil
}

// Save saves the backups to yaml file.
func (backups *Backups) Save() error {
	backupFile := config.Conf.General.BackupFile
	f, err := os.Create(backupFile)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("Failed to save backup file. err: %s", err)
	}
	return yaml.NewEncoder(f).Encode(backups)
}

func (backups *Backups) ToString() (string, error) {
	var buffer bytes.Buffer
	err := yaml.NewEncoder(&buffer).Encode(backups)
	if err != nil {
		return "", fmt.Errorf("Failed to convert struct to yaml string: %v", err)
	}
	return buffer.String(), nil
}

func (backups *Backups) Order() {
	sortBy := config.Conf.General.SortBy
	switch {

	case sortBy == "-recency":
		backups.Reverse()
	}
}

func (backups *Backups) Reverse() {
	for i, j := 0, len(backups.Files)-1; i < j; i, j = i+1, j-1 {
		backups.Files[i], backups.Files[j] = backups.Files[j], backups.Files[i]
	}
}
