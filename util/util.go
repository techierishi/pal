package util

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"github.com/ulfox/dby/db"
)

var (
	PalShellrc  = "[ -f ~/.palrc ] && source ~/.palrc # added by PalApp"
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func EditFile(command, file string) error {
	command += " " + file
	return RunCmd(command, os.Stdin, os.Stdout)
}

func runPowerShellCommandInBackground(command string) (int, error) {
	psCmd := exec.Command("powershell", "-Command", command)

	err := psCmd.Start()
	if err != nil {
		return 0, fmt.Errorf("error starting PowerShell command: %v", err)
	}

	return psCmd.Process.Pid, nil
}

func RunCmdInBackground(commandToRun string) int {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		if isPowerShell() {
			pid, err := runPowerShellCommandInBackground(commandToRun)
			if err != nil {
				fmt.Printf("error to run cmd in powershell %v \n", err)
			}
			return pid
		}
		cmd = exec.Command("cmd", "/C", "start", "cmd", "/b", commandToRun)
	} else {
		cmd = exec.Command("/bin/sh", "-c", commandToRun, "&")

	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		os.Exit(1)
	}

	fmt.Printf("Pal service started with PID: %d\n", cmd.Process.Pid)
	return cmd.Process.Pid
}

func RunCmd(command string, r io.Reader, w io.Writer) error {
	var cmd *exec.Cmd
	if len(config.Conf.General.Cmd) > 0 {
		line := append(config.Conf.General.Cmd, command)
		cmd = exec.Command(line[0], line[1:]...)
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Stderr = os.Stderr
	cmd.Stdout = w
	cmd.Stdin = r
	return cmd.Run()
}

func RandStrRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ReadableTime() string {
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02_15:04:05")

	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	timeString = re.ReplaceAllString(timeString, "_")

	return timeString
}
func UnixMilli() int64 {
	return time.Now().UnixMilli()
}

func UnixMilliDiff(milliseconds int64) int {
	timeFromMilliseconds := time.Unix(0, milliseconds*int64(time.Millisecond))
	duration := time.Since(timeFromMilliseconds)
	return int(duration)
}
func isPowerShell() bool {
	psModulePath := os.Getenv("PSModulePath")
	return psModulePath != ""
}

func GetCurrentShell() string {
	if runtime.GOOS == "windows" {
		if isPowerShell() {
			return "powershell"
		}
		if os.Getenv("ComSpec") != "" {
			return "cmd"
		}
	}
	shell := strings.ToLower(os.Getenv("SHELL"))
	if strings.Contains(shell, "bash") {
		return "bash"
	} else if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	return "unknown"
}

func ParseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "Yes", "YES", "yes", "Y", "y":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "No", "NO", "no", "N", "n":
		return false, nil
	}
	return false, errors.New("ParseBool error")
}

func initCheckFromCache(stateFile *db.Storage) bool {

	lastCheckTime, err := stateFile.GetPath("init.lastStatusTime")
	if err != nil {
		return false
	}
	duration := UnixMilliDiff(int64(lastCheckTime.(int)))
	if duration >= 24 {
		return false
	}

	initCheckBool, err := stateFile.GetPath("app.status")
	if err != nil {
		return false
	}
	return initCheckBool.(bool)
}

func CheckIfInitRan() error {
	logger := logr.GetLogInstance()

	if len(os.Args) >= 2 && os.Args[1] == "init" {
		// Skip the check if its init itself
		return nil
	}

	stateFile, err := GetStateFile(logger)
	if err != nil {
		logger.Error().Any("error", err).Msg("Error opening app db")
	}
	defer stateFile.Close()

	if initCheckFromCache(stateFile) {
		return nil
	}

	shell := GetCurrentShell()

	if strings.EqualFold(shell, "cmd") || strings.EqualFold(shell, "powershell") {
		// Init only supported for zsh and bash
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		return err
	}
	shellRcPath := homeDir + "/.bashrc"

	if strings.EqualFold(shell, "zsh") {
		shellRcPath = homeDir + "/.zshrc"
	}

	contains, err := FileContains(shellRcPath, PalShellrc)
	if err != nil {
		fmt.Println("Error checking "+shellRcPath, err)
		return err
	}

	if !contains {
		fmt.Fprintf(color.Output, "%12s", color.YellowString("Ensure you've run the `pal init` \n"))
		return fmt.Errorf("Palrc not generated")
	}

	stateFile.Upsert("init.status", true)
	stateFile.Upsert("init.lastStatusTime", UnixMilli())

	return nil

}
