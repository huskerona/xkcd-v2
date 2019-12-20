package file_manager

import (
	"encoding/gob"
	"fmt"
	"github.com/huskerona/xkcd2/infrastructure/model"
	"github.com/huskerona/xkcd2/infrastructure/util"
	"io"
	"os"
)

// Writes comics into an index file. This process will recreate the file every time.
// Better approach would be to find what has been written before and append the new items.
// (Will be done later)
func WriteIndexFile(comics []*model.XKCD) error {
	file, err := os.OpenFile(util.GetIndexFile(), os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		return fmt.Errorf("gob WriteIndexFile: %v", err)
	}

	defer file.Close()

	err = gob.NewEncoder(file).Encode(&comics)

	if err != nil {
		return fmt.Errorf("gob encode: %v", err)
	}

	return nil
}

// Reads the index file and loads all the comics into a slice.
func ReadIndexFile() ([]*model.XKCD, error) {
	file, err := os.OpenFile(util.GetIndexFile(), os.O_RDONLY, 0777)

	if err != nil {
		return nil, fmt.Errorf("gob ReadIndexFile: %v", err)
	}

	defer file.Close()

	var comics []*model.XKCD
	var current *model.XKCD

	decoder := gob.NewDecoder(file)

	for decoder.Decode(current) != io.EOF {
		comics = append(comics, current)
	}

	return comics, nil
}
