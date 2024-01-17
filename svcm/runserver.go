package svcm

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
)

type Server struct {
	logger *zerolog.Logger
}

func NewServer(logger *zerolog.Logger) Server {
	return Server{
		logger: logger,
	}
}
func (s *Server) Run() {
	var signalSubscribers []chan os.Signal

	shutdown := make(chan string)

	// handlers
	mux := http.NewServeMux()
	mux.Handle("/status", &statusHandler{})

	server := &http.Server{
		Addr:              "localhost:" + strconv.Itoa(config.Conf.Daemon.Port),
		Handler:           mux,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	go server.ListenAndServe()

	SigRun(s.logger, signalSubscribers, shutdown, server)
}
