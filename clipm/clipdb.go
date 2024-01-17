package clipm

import (
	"sync"

	"github.com/rs/zerolog"
	"github.com/techierishi/pal/config"
	"github.com/ulfox/dby/db"
)

var clipDb *db.Storage
var clipDbOnce sync.Once

func GetClipDb(logger *zerolog.Logger) (*db.Storage, error) {
	clipDbOnce.Do(func() {
		var err error

		clipDb, err = db.NewStorageFactory(config.Conf.General.ClipboardFile)
		if err != nil {
			logger.Fatal().AnErr("Error opening clip db", err)
		}
	})
	return clipDb, nil
}
