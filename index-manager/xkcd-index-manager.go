package index_manager

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	imageSvc "github.com/huskerona/xkcd2/image-service"
	"github.com/huskerona/xkcd2/infrastructure"
	"github.com/huskerona/xkcd2/infrastructure/logger"
	"github.com/huskerona/xkcd2/infrastructure/model"
	netManager "github.com/huskerona/xkcd2/net-manager"
)

var mu sync.Mutex

//+ Exports

func GetComics() []*model.XKCD {
	defer logger.Trace("func GetComics")()

	return model.Comics
}

// Returns a comic based on the comic number, or nil if comic was not found.
func GetComic(comicNum int) *model.XKCD {
	defer logger.Trace(fmt.Sprintf("func GetComic(%d)", comicNum))()

	index, comic := model.Comics.Get(comicNum)

	if index == -1 {
		return nil
	}

	return comic
}

func LoadComics(comics []*model.XKCD) {
	defer logger.Trace("func LoadComics")()

	mu.Lock()
	defer mu.Unlock()
	for _, item := range comics {
		model.Comics = append(model.Comics, item)
	}

}

// Downloads the XKCD comic based on its number. If number is 0 it will
// fetch the latest issue.
func DownloadComic(comicNumber int) (*model.XKCD, error) {
	defer logger.Trace("func DownloadComic")()

	var err error
	var xkcd *model.XKCD

	if comicNumber == 0 {
		xkcd, err = fetchLatestComic()
	} else {
		xkcd, err = fetchComic(comicNumber)
	}

	if err != nil {
		return nil, err
	}

	imageByte, err := downloadImage((*xkcd).ImageURL)

	if err != nil {
		// There are some comics whose image cannot be retrieved. It would require
		// that we parse the HTML. In any case, XKCD struct is returned.
		return xkcd, nil
	}

	(*xkcd).Image = imageSvc.EncodeToBase64(imageByte)

	return xkcd, nil
}

// Adds a comic to the collection
func AddToCollection(xkcd *model.XKCD) {
	defer logger.Trace(fmt.Sprintf("func AddToCollection(%d)", (*xkcd).Number))()

	model.Comics.Add(xkcd)
}

// Returns true if the comic exists in the offline collection or false if not.
func Contains(comicNumber int) bool {
	defer logger.Trace("func Contains()")()

	return model.Comics.Contains(comicNumber)
}

// Returns the total number of items in the index
func Count() int {
	defer logger.Trace("func Count()")()

	return len(model.Comics)
}

func Sort() {
	mu.Lock()
	defer mu.Unlock()
	sort.Sort(model.Comics)
}

//- Exports

//+ Helper funcs

// Retrieves the latest comic and returns the XKCD
func fetchLatestComic() (*model.XKCD, error) {
	defer logger.Trace("func fetchLatestComic()")()

	url := fmt.Sprintf("%s/%s", infrastructure.HomeURL, infrastructure.JSONURL)

	rawJson, err := netManager.GetDataFromURL(url)

	if err != nil {
		return nil, err
	}

	// Read JSON from body as string
	var xkcd *model.XKCD
	if err := json.Unmarshal(rawJson, &xkcd); err != nil {
		return nil, fmt.Errorf("unmarshall comic json: %v", err)
	}

	return xkcd, nil
}

// Retrieves a specific comic and returns the XKCD
func fetchComic(comicNumber int) (*model.XKCD, error) {
	defer logger.Trace("fetch fetchComic()")()

	url := fmt.Sprintf("%s/%d/%s", infrastructure.HomeURL, comicNumber, infrastructure.JSONURL)

	rawJson, err := netManager.GetDataFromURL(url)

	if err != nil {
		return nil, err
	}

	var xkcd *model.XKCD
	if err := json.Unmarshal(rawJson, &xkcd); err != nil {
		return nil, fmt.Errorf("unmarshall comic json: %v", err)
	}

	return xkcd, nil
}

// Downloads an image found in the ImageURL field in XKCD struct.
func downloadImage(url string) ([]byte, error) {
	defer logger.Trace("func downloadImage()")()

	return netManager.GetDataFromURL(url)
}

//- Helper funcs
