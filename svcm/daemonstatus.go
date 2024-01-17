package svcm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type Daemon struct {
	Logger *zerolog.Logger
}

func (d *Daemon) KillDaemon(pidStr string) error {

	pid, err := strconv.Atoi(strings.TrimSuffix(pidStr, "\n"))
	if err != nil {
		return fmt.Errorf("could not parse PID file contents: %w", err)
	}
	err = StopProcess(pid)
	if err != nil {
		return fmt.Errorf("could not stop pal process with PID %d: %w", pid, err)
	}
	return nil
}

func StopProcess(pid int) error {
	if runtime.GOOS == "windows" {
		err := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid)).Run()
		if err != nil {
			return fmt.Errorf("kill command finished with error: %w", err)
		}
	} else {
		err := exec.Command("kill", "-SIGTERM", fmt.Sprintf("%d", pid)).Run()
		if err != nil {
			return fmt.Errorf("kill command finished with error: %w", err)
		}
	}

	return nil
}

func get(port int) (*http.Response, error) {
	url := "http://localhost:" + strconv.Itoa(port) + "/status"
	client := NewHttpClient()
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error while GET'ing daemon /status: %w", err)
	}
	return resp, nil
}

func IsDaemonRunning(port int) (bool, error) {
	resp, err := get(port)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return true, nil
}

func GetDaemonStatus(port int) (*StatusResponse, error) {
	resp, err := get(port)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	jsn, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading 'daemon /status' response: %w", err)
	}
	var msgResp StatusResponse
	err = json.Unmarshal(jsn, &msgResp)
	if err != nil {
		return nil, fmt.Errorf("error while decoding 'daemon /status' response: %w", err)
	}
	return &msgResp, nil
}
