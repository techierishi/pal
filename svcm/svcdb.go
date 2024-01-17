package svcm

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
	"github.com/ulfox/dby/db"
)

var svcDb *db.Storage
var svcDbOnce sync.Once

func GetSvcDb(logger *zerolog.Logger) (*db.Storage, error) {
	svcDbOnce.Do(func() {
		var err error
		svcDb, err = db.NewStorageFactory(config.Conf.General.SvcFile)
		if err != nil {
			logger.Fatal().AnErr("Error opening svc db", err)
		}
	})
	return svcDb, nil
}
