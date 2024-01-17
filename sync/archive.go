package sync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/techierishi/pal/config"
	"gopkg.in/yaml.v3"
)

func RestoreFiles(bkpFilePath string, restoreDir string) error {

	if _, err := os.Stat(bkpFilePath); err != nil {
		fmt.Println("Restore file path does not exist.")
		return err
	}

	file, err := os.Open(bkpFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

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

	for {
		var data FileInfo
		if err := decoder.Decode(&data); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error decoding data:", err)
			return err
		}

		if _, ok := data.Metadata["name"]; !ok {
			continue
		}
		var fileWritePath string

		dirPath := data.Metadata["path"]
		if !strings.HasPrefix(dirPath, "<home>") && !strings.HasPrefix(dirPath, "<config>") {
			fmt.Println("Only filepath from home or config folder is supported for restore. Skipping... ", dirPath)
			continue
		}
		dirPath = strings.Replace(dirPath, "~", homeDir, 1)

		dirPath = strings.Replace(dirPath, "<config>", confDir, 1)
		dirPath = strings.Replace(dirPath, "<home>", homeDir, 1)

		fileWritePath = filepath.Join(dirPath, data.Metadata["name"])
		err = os.WriteFile(fileWritePath, []byte(data.FileContent), os.ModePerm)

		fmt.Printf("Restored file: %s\n", data.Metadata["name"])
	}
	return nil
}

func BackupFiles(filesPaths []string, outputPath string) error {

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Saving backup..."
	s.Start()

	file, err := os.OpenFile(config.Conf.General.BackupFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error creating backup file:", err)
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)

	if err := encoder.Encode(FileInfo{
		Metadata: map[string]string{
			"updated_at": time.Now().Format(time.RFC3339),
		},
	}); err != nil {
		fmt.Println("Error encoding data:", err)
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("BackupFiles:: Error getting user home directory:", err)
		return err
	}

	confDir, err := config.GetUserConfigDir()
	if err != nil {
		fmt.Println("BackupFiles:: Error getting user config directory:", err)
		return err
	}

	for _, filePath := range filesPaths {

		content, err := readFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}

		dirPath := filepath.Dir(filePath)
		if !strings.HasPrefix(dirPath, homeDir) || !strings.HasPrefix(dirPath, confDir) {
			fmt.Println("Only filepath from home or config folder is supported for backup. Skipping... ", dirPath)
			continue
		}

		dirPath = strings.Replace(dirPath, confDir, "<config>", 1)
		dirPath = strings.Replace(dirPath, homeDir, "<home>", 1)

		fileData := FileInfo{
			Metadata: map[string]string{
				"path": dirPath,
				"name": filepath.Base(filePath),
			},
			FileContent: string(content),
		}
		if err := encoder.Encode(fileData); err != nil {
			fmt.Println("Error encoding data:", err)
			return err
		}
	}

	if err := copyFile(config.Conf.General.BackupFile, filepath.Base(config.Conf.General.BackupFile)); err != nil {
		fmt.Println("Error creating backup:", err)
		return err
	}

	s.Stop()
	fmt.Fprintf(color.Output, "%12s", color.GreenString(fmt.Sprintf("Backup file created successfully. \n")))
	return nil
}

func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func copyFile(sourceFilePath, destinationFileName string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFilePath := filepath.Join(".", destinationFileName)
	destinationFile, err := os.Create(destinationFilePath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

func writeFile(filePath string, data []byte, perm os.FileMode) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return file.Chmod(perm)
}
