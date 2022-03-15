package index_manager

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	imageSvc "xkcd2/image-service"
	"xkcd2/infrastructure"
	"xkcd2/infrastructure/logger"
	"xkcd2/infrastructure/model"
	netManager "xkcd2/net-manager"
)

var mu sync.Mutex

//+ Exports

func GetComics() []*model.XKCD {
	defer logger.Trace("func GetComics")()

	return model.Comics
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

	imageByte, err := downloadImage(xkcd.ImageURL) // prev: (*xkcd).ImageURL

	if err != nil {
		// There are some comics whose image cannot be retrieved. It would require
		// that we parse the HTML. In any case, XKCD struct is returned.
		return xkcd, nil
	}

	xkcd.Image = imageSvc.EncodeToBase64(imageByte)

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
	xkcd, err := fetch(url)

	if err != nil {
		return nil, fmt.Errorf("fetchLatestComic: %v", err)
	}

	return xkcd, nil
}

// Retrieves a specific comic and returns the XKCD
func fetchComic(comicNumber int) (*model.XKCD, error) {
	defer logger.Trace("fetch fetchComic()")()

	url := fmt.Sprintf("%s/%d/%s", infrastructure.HomeURL, comicNumber, infrastructure.JSONURL)
	xkcd, err := fetch(url)

	if err != nil {
		return nil, fmt.Errorf("fetchComic: %v", err)
	}

	return xkcd, nil
}

// downloadImage fetches an image specified in imageUrl parameter.
func downloadImage(imageUrl string) ([]byte, error) {
	defer logger.Trace("func downloadImage()")()

	result, err := netManager.Get(imageUrl)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func fetch(url string) (*model.XKCD, error) {
	result, err := netManager.Get(url)

	if err != nil {
		return nil, err
	}

	var xkcd *model.XKCD
	if err := json.Unmarshal(result, &xkcd); err != nil {
		return nil, fmt.Errorf("fetch: %v", err)
	}

	return xkcd, nil
}

//- Helper funcs
