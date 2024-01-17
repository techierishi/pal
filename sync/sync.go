package sync

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/techierishi/pal/config"
)

// Client manages communication with the remote Snippet repository
type Client interface {
	GetSnippet() (*Snippet, error)
	UploadSnippet(string) error
}

// Snippet is the remote snippet
type Snippet struct {
	Content   string
	UpdatedAt time.Time
}

// AutoSync syncs snippets automatically
func AutoSync(file string) error {
	client, err := NewSyncClient()
	if err != nil {
		return errors.Wrap(err, "Failed to initialize API client")
	}

	snippet, err := client.GetSnippet()
	if err != nil {
		return err
	}

	fi, err := os.Stat(file)
	if os.IsNotExist(err) || fi.Size() == 0 {
		return download(snippet.Content)
	} else if err != nil {
		return errors.Wrap(err, "Failed to get a FileInfo")
	}

	local := fi.ModTime().UTC()
	remote := snippet.UpdatedAt.UTC()

	switch {
	case local.After(remote):
		return upload(client)
	case remote.After(local):
		return download(snippet.Content)
	default:
		return nil
	}
}

// NewSyncClient returns Client
func NewSyncClient() (Client, error) {
	if config.Conf.General.Backend == "gitlab" {
		client, err := NewGitLabClient()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to initialize GitLab client")
		}
		return client, nil
	}
	client, err := NewGistClient()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize Gist client")
	}
	return client, nil
}

func upload(client Client) (err error) {
	body, err := readFile(config.Conf.General.BackupFile)
	if err != nil {
		return err
	}

	if err = client.UploadSnippet(string(body)); err != nil {
		return errors.Wrap(err, "Failed to upload snippet")
	}

	fmt.Println("Upload success")
	return nil
}

func download(content string) error {
	backupFile := config.Conf.General.BackupFile

	body, err := readFile(config.Conf.General.BackupFile)
	if err != nil {
		return err
	}
	if content == string(body) {
		fmt.Println("Already up-to-date")
		return nil
	}

	fmt.Println("Download success")
	return os.WriteFile(backupFile, []byte(content), os.ModePerm)
}
