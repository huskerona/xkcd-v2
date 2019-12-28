package shared

import (
	fileManager "github.com/huskerona/xkcd2/file-manager"
	"github.com/huskerona/xkcd2/infrastructure/logger"
	"github.com/huskerona/xkcd2/infrastructure/model"
	"log"
)

// Loads the comics from the index file
func LoadComics() ([]*model.XKCD, error) {
	defer logger.Trace("loadComics")()

	comics, err := fileManager.ReadIndexFile()

	return comics, err
}


// Writes the comics back to the index file
func WriteComics(comics []*model.XKCD) {
	err := fileManager.WriteIndexFile(comics)

	if err != nil {
		log.Fatal(err)
	}
}
