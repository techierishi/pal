package logr

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
)

var logger *zerolog.Logger
var once sync.Once

func GetLogInstance() *zerolog.Logger {
	once.Do(func() {
		var err error
		logger, err = NewLogger(zerolog.DebugLevel)
		if err != nil {
			log.Fatal("error getting logger instance ", err)
		}
	})
	return logger
}

func NewLogger(level zerolog.Level) (*zerolog.Logger, error) {
	dataDir, err := config.GetDefaultConfigDir()
	if err != nil {
		return nil, fmt.Errorf("error while getting RESH data dir: %w", err)
	}
	logPath := filepath.Join(dataDir, "pal.log")

	file, err := os.OpenFile(
		logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	// defer file.Close()

	logger := zerolog.New(file).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(level)

	return &logger, nil
}
