package svcm

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

func sendSignals(logger *zerolog.Logger, sig os.Signal, subscribers []chan os.Signal, done chan string) {
	for _, sub := range subscribers {
		sub <- sig
	}
	logger.Warn().Msg("Sent shutdown signals to components")
	chanCount := len(subscribers)
	start := time.Now()
	delay := time.Millisecond * 100
	timeout := time.Millisecond * 2000

	for {
		select {
		case _ = <-done:
			chanCount--
			if chanCount == 0 {
				logger.Warn().Msg("All components shut down successfully")
				return
			}
		default:
			time.Sleep(delay)
		}
		if time.Since(start) > timeout {
			logger.Error().
				Str("componentsStillUp", strconv.Itoa(chanCount)).
				Str("timeout", timeout.String()).
				Msg("Timeouted while waiting for proper shutdown")
			return
		}
	}
}

// Run catches and handles signals
func SigRun(logger *zerolog.Logger, subscribers []chan os.Signal, done chan string, server *http.Server) {
	signals := make(chan os.Signal, 1)
	loggerSigMod := logger.With().Str("module", "signalhandler").Logger()

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	var sig os.Signal
	for {
		sig := <-signals
		loggerSig := loggerSigMod.With().Str("signal", sig.String()).Logger()
		loggerSig.Info().Msg("Got signal")
		if sig == syscall.SIGTERM {
			// Shutdown daemon on SIGTERM
			break
		}
		loggerSig.Warn().Msg("Ignoring signal. Send SIGTERM to trigger shutdown.")
	}

	loggerSigMod.Info().Msg("Sending shutdown signals to components ...")
	sendSignals(logger, sig, subscribers, done)

	loggerSigMod.Info().Msg("Shutting down the server ...")
	if err := server.Shutdown(context.Background()); err != nil {
		loggerSigMod.Error().AnErr("Error while shuting down HTTP server", err)
	}
}
