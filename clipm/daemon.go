package clipm

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/techierishi/pal/logr"
	"github.com/techierishi/pal/util"
	"golang.design/x/clipboard"
)

type Clip struct {
	ID      int
	Time    int64
	Content []byte
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func Record() error {
	logger := logr.GetLogInstance()
	logger.Info().Msg("Clipboard recording started...")

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	clipDb, err := GetClipDb(logger)
	defer clipDb.Close()
	ch := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for data := range ch {

		copiedStr := string(data)

		timestamp := util.UnixMilli()
		clipInfo := ClipInfo{
			Timestamp: timestamp,
			Content:   copiedStr,
		}
		hash := util.CalculateHash(copiedStr)

		clipDb.Upsert(fmt.Sprintf("%s.%s", CLIPBOARD_TBL, hash), clipInfo)
		str := util.CleanStr(copiedStr).StandardizeSpaces().TruncateText(10).ReplaceNewLine()
		logger.Info().Msg(string(str + "... COPIED!"))

	}

	return nil
}
