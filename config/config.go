package config

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Conf is global config variable
var Conf Config

// Config is a struct of config
type Config struct {
	General GeneralConfig `yaml:"General"`
	Gist    GistConfig    `yaml:"Gist"`
	GitLab  GitLabConfig  `yaml:"GitLab"`
	Daemon  DaemonConfig  `yaml:"Daemon"`
}

// GeneralConfig is a struct of general config
type GeneralConfig struct {
	SnippetFile   string `yaml:"snippetfile"`
	CredFile      string `yaml:"credfile"`
	AliasFile     string `yaml:"aliasfile"`
	BackupFile    string `yaml:"backupfile"`
	ClipboardFile string `yaml:"clipboardfile"`
	SvcFile       string `yaml:"svcfile"`
	// General purpose app state
	StateFile string   `yaml:"statefile"`
	Editor    string   `yaml:"editor"`
	Column    int      `yaml:"column"`
	Backend   string   `yaml:"backend"`
	SortBy    string   `yaml:"sortby"`
	Cmd       []string `yaml:"cmd"`
}

type DaemonConfig struct {
	Port int `yaml:"port"`
}

// GistConfig is a struct of config for Gist
type GistConfig struct {
	FileName    string `yaml:"file_name"`
	AccessToken string `yaml:"access_token"`
	GistID      string `yaml:"gist_id"`
	Public      bool   `yaml:"public"`
	AutoSync    bool   `yaml:"auto_sync"`
}

// GitLabConfig is a struct of config for GitLabSnippet
type GitLabConfig struct {
	FileName    string `yaml:"file_name"`
	AccessToken string `yaml:"access_token"`
	Url         string `yaml:"url"`
	ID          string `yaml:"id"`
	Visibility  string `yaml:"visibility"`
	AutoSync    bool   `yaml:"auto_sync"`
	Insecure    bool   `yaml:"skip_ssl"`
}

// Flag is global flag variable
var Flag FlagConfig

// FlagConfig is a struct of flag
type FlagConfig struct {
	Debug        bool
	Query        string
	FilterTag    string
	Command      bool
	Copy         bool
	Delimiter    string
	OneLine      bool
	Color        bool
	Tag          bool
	Detach       bool
	HasClipboard bool
}

// Load loads a config yaml
func (cfg *Config) Load(file string) error {
	_, err := os.Stat(file)
	if err == nil {

		snippetFile, err := os.OpenFile(file, os.O_RDONLY, 0600)
		if err != nil {
			return err
		}
		defer snippetFile.Close()
		dec := yaml.NewDecoder(snippetFile)
		err = dec.Decode(&cfg)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Config file is empty.")
				return nil
			}
			return fmt.Errorf("Failed to load config file. %v", err)
		}

		cfg.General.SnippetFile = expandPath(cfg.General.SnippetFile)
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	dir, err := GetDefaultConfigDir()
	if err != nil {
		return errors.Wrap(err, "Failed to get the default config directory")
	}
	cfg.Daemon.Port = 7200
	cfg.General.SnippetFile = filepath.Join(dir, "snippet.yaml")
	cfg.General.CredFile = filepath.Join(dir, "credential.yaml")
	cfg.General.AliasFile = filepath.Join(dir, "alias.yaml")
	cfg.General.BackupFile = filepath.Join(dir, "pal-backups.yaml")
	cfg.General.ClipboardFile = filepath.Join(dir, "clipboard.yaml")
	cfg.General.StateFile = filepath.Join(dir, "state.yaml")
	cfg.General.SvcFile = filepath.Join(dir, "svc.yaml")

	cfg.General.Editor = os.Getenv("EDITOR")
	if cfg.General.Editor == "" && runtime.GOOS != "windows" {
		if isCommandAvailable("sensible-editor") {
			cfg.General.Editor = "sensible-editor"
		} else {
			cfg.General.Editor = "vim"
		}
	}
	cfg.General.Column = 40
	cfg.General.Backend = "gist"

	cfg.Gist.FileName = "pal-backups.yaml"

	cfg.GitLab.FileName = "pal-backups.yaml"
	cfg.GitLab.Visibility = "private"

	return yaml.NewEncoder(f).Encode(cfg)
}

func GetUserConfigDir() (dir string, err error) {
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "pal")
		}
		dir = filepath.Join(dir, "pal")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "pal")
	}

	return dir, err
}

// GetDefaultConfigDir returns the default config directory
func GetDefaultConfigDir() (dir string, err error) {
	if env, ok := os.LookupEnv("PAL_CONFIG_DIR"); ok {
		dir = env
	} else if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "pal")
		}
		dir = filepath.Join(dir, "pal")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "pal")
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.Join(os.Getenv("USERPROFILE"), s[2:])
		} else {
			s = filepath.Join(os.Getenv("HOME"), s[2:])
		}
	}
	return os.Expand(s, os.Getenv)
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
