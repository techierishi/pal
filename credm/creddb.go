package credm

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
	"github.com/ulfox/dby/db"
)

var credDb *db.Storage
var credDbOnce sync.Once

func GetCredDb(logger *zerolog.Logger) (*db.Storage, error) {
	credDbOnce.Do(func() {
		var err error

		credDb, err = db.NewStorageFactory(config.Conf.General.CredFile)
		if err != nil {
			logger.Fatal().AnErr("Error opening cred db", err)
		}
	})
	return credDb, nil
}
