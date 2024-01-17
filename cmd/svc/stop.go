package svc

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/svcm"
)

var clipStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop service",
	RunE:  stopFunc,
}

func init() {
	RootCmdSvc.AddCommand(clipStopCmd)
}

func stopFunc(cmd *cobra.Command, args []string) (err error) {
	logger := logr.GetLogInstance()
	d := svcm.Daemon{Logger: logger}
	pidLogger := logger.With().Int("daemonPID", os.Getpid()).Logger()

	svcDb, err := svcm.GetSvcDb(logger)
	if err != nil {
		pidLogger.Error().Any("error", err).Msg("Error opening app db")
	}
	defer svcDb.Close()
	res, err := svcm.IsDaemonRunning(config.Conf.General.Column)

	pidLogger.Error().Any("error", err).Msg("Error while checking daemon status - it's probably not running")

	if res {
		pidLogger.Error().Msg("Daemon is already running - exiting!")
		return
	}
	pidStr, err := svcDb.GetPath("pid")
	if err != nil {
		pidLogger.Error().Any("error", err).Msg("Could not kill daemon")
	}

	err = d.KillDaemon(pidStr.(string))
	if err != nil {
		pidLogger.Error().Any("error", err).Msg("Could not kill daemon")
	}

	return nil
}
