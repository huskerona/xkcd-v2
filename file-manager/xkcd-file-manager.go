package file_manager

import (
	"encoding/gob"
	"fmt"
	"github.com/huskerona/xkcd2/infrastructure/logger"
	"github.com/huskerona/xkcd2/infrastructure/model"
	"github.com/huskerona/xkcd2/infrastructure/util"
	"os"
)

func init() {
	// Probably not the best place as init might fail. Will change.
	if _, err := os.Stat(util.GetXkcdFolder()); err != nil {
		os.Mkdir(util.GetXkcdFolder(), 0777)
	}
}
// Writes comics into an index file. This process will recreate the file every time.
// Better approach would be to find what has been written before and append the new items.
// (Will be done later)
func WriteIndexFile(comics []*model.XKCD) error {
	defer logger.Trace("WriteIndexFile")()

	file, err := os.OpenFile(util.GetIndexFile(), os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		return fmt.Errorf("gob WriteIndexFile: %v", err)
	}

	defer file.Close()

	encoder := gob.NewEncoder(file)

	for _, current := range comics {
		err = encoder.Encode(&current)

		if err != nil {
			return fmt.Errorf("gob encode: %v", err)
		}

		logger.Info(fmt.Sprintf("Encoding %d (%s-%s-%s)\n",
			current.Number, current.Year, current.Month, current.Day))
	}

	return nil
}

// Reads the index file and loads all the comics into a slice.
func ReadIndexFile() ([]*model.XKCD, error) {
	defer logger.Trace("ReadingIndexFile")()

	file, err := os.OpenFile(util.GetIndexFile(), os.O_RDONLY, 0777)

	if err != nil {
		return nil, fmt.Errorf("gob ReadIndexFile: %v", err)
	}

	defer file.Close()

	logger.Info(fmt.Sprintf("ReadIndexFile file opened, decoding\n"))

	loaded := false
	var comics []*model.XKCD

	decoder := gob.NewDecoder(file)

	for !loaded {
		current := &model.XKCD{}

		if err := decoder.Decode(&current); err == nil {
			comics = append(comics, current)

			logger.Info(fmt.Sprintf("Decoding %d (%s-%s-%s)\n",
				(*current).Number, (*current).Year, (*current).Month, (*current).Day))
		} else {
			loaded = true
		}
	}

	logger.Info(fmt.Sprintf("ReadIndexFile completed with total of %d\n", len(comics)))

	return comics, nil
}
