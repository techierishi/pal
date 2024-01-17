package util

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
	"github.com/ulfox/dby/db"
)

var sateDb *db.Storage
var sateDbOnce sync.Once

func GetStateFile(logger *zerolog.Logger) (*db.Storage, error) {
	sateDbOnce.Do(func() {
		var err error
		sateDb, err = db.NewStorageFactory(config.Conf.General.StateFile)
		if err != nil {
			logger.Fatal().AnErr("Error opening sate db", err)
		}
	})
	return sateDb, nil
}
