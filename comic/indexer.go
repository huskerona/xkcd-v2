package comic

import (
	"fmt"
	"sort"

	"xkcd2/tools/logger"
)

func GetComics() []*XKCD {
	defer logger.Trace("func GetComics")()

	return Comics
}

func LoadComics(comics []*XKCD) {
	defer logger.Trace("func LoadComics")()

	mu.Lock()
	defer mu.Unlock()
	for _, item := range comics {
		Comics = append(Comics, item)
	}
}

// Adds a comic to the collection
func AddToCollection(xkcd *XKCD) {
	defer logger.Trace(fmt.Sprintf("func AddToCollection(%d)", (*xkcd).Number))()

	Comics.Add(xkcd)
}

// Returns true if the comic exists in the offline collection or false if not.
func Contains(comicNumber int) bool {
	defer logger.Trace("func Contains()")()

	return Comics.Contains(comicNumber)
}

// Returns the total number of items in the index
func Count() int {
	defer logger.Trace("func Count()")()

	return len(Comics)
}

func Sort() {
	mu.Lock()
	defer mu.Unlock()
	sort.Sort(Comics)
}
