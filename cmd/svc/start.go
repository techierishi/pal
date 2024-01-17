package svc

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/techierishi/pal/clipm"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/svcm"
	"github.com/techierishi/pal/util"
)

// info passed during build
var commit string
var development string

const helpMsg = `ERROR: pal daemon doesn't accept any arguments`

var cliStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start service",
	RunE:  startFunc,
}

func init() {
	RootCmdSvc.AddCommand(cliStartCmd)
	RootCmdSvc.PersistentFlags().BoolVarP(&config.Flag.Detach, "detach", "d", false, "detach mode")

}
func startFunc(cmd *cobra.Command, args []string) (err error) {

	if config.Flag.Detach {
		return startDetach()
	}

	logger := logr.GetLogInstance()
	svcDb, err := svcm.GetSvcDb(logger)
	if err != nil {
		logger.Error().Any("error", err).Msg("Error opening app db")
	}

	defer svcDb.Close()
	d := svcm.Daemon{Logger: logger}
	logger.Info().
		Str("version", version).
		Str("commit", commit).Msg("Daemon starting ...")

	pidLogger := logger.With().Int("daemonPID", os.Getpid()).Logger()

	res, err := svcm.IsDaemonRunning(config.Conf.Daemon.Port)
	if err != nil {
		pidLogger.Error().Any("error", err).Msg("Error while checking daemon status - it's probably not running")
	}
	if res {
		pidLogger.Error().Msg("Daemon is already running - exiting!")
		return
	}

	// Kill and Delete the process if it exists
	pidStr, err := svcDb.GetPath("pid")
	if err == nil {
		err = d.KillDaemon(pidStr.(string))
		if err != nil {
			pidLogger.Error().Any("error", err).Msg("Could not kill daemon")
		}
	}

	err = svcDb.Upsert("pid", strconv.Itoa(os.Getpid()))
	if err != nil {
		pidLogger.Error().Any("error", err).Msg("Could not save pid")
	}
	defer svcDb.Delete("pid")
	go clipm.Record()
	server := svcm.NewServer(logger)
	server.Run()
	logger.Info().Msg("Shutting down ...")

	return nil
}

func startDetach() (err error) {
	logger := logr.GetLogInstance()

	pidLogger := logger.With().Int("daemonPID", os.Getpid()).Logger()

	res, err := svcm.IsDaemonRunning(config.Conf.Daemon.Port)
	if err != nil {
		pidLogger.Error().Any("error", err).Msg("Error while checking daemon status - it's probably not running")
	}
	if res {
		pidLogger.Error().Msg("Daemon is already running - exiting!")
		return nil
	}
	util.RunCmdInBackground("pal svc start")
	return nil
}
