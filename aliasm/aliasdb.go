package aliasm

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
	"github.com/ulfox/dby/db"
)

var aliasDb *db.Storage
var aliasDbOnce sync.Once

func GetAliasDb(logger *zerolog.Logger) (*db.Storage, error) {
	aliasDbOnce.Do(func() {
		var err error
		aliasDb, err = db.NewStorageFactory(config.Conf.General.AliasFile)
		if err != nil {
			logger.Fatal().AnErr("Error opening alias db", err)
		}
	})
	return aliasDb, nil
}
