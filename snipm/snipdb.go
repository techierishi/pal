package snipm

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
	"github.com/ulfox/dby/db"
)

var snipDb *db.Storage
var snipDbOnce sync.Once

func GetSnipDb(logger *zerolog.Logger) (*db.Storage, error) {
	snipDbOnce.Do(func() {
		var err error
		snipDb, err = db.NewStorageFactory(config.Conf.General.SnippetFile)
		if err != nil {
			logger.Fatal().AnErr("Error opening snip db", err)
		}
	})
	return snipDb, nil
}
