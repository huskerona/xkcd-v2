package comic

import (
	"encoding/json"
	"fmt"
	"xkcd2/tools/logger"
	"xkcd2/webclient"
)

// downloadImage fetches an image specified in imageUrl parameter.
func downloadImage(imageUrl string) ([]byte, error) {
	defer logger.Trace(fmt.Sprintf("func downloadImage(%s)", imageUrl))()

	result, err := webclient.Get(imageUrl)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// fetch perfroms GET operation on url and it unmarshalls the JSON document into XKCD object.
func fetch(url string) (*XKCD, error) {
	defer logger.Trace(fmt.Sprintf("func fetch(%s)", url))()
	result, err := webclient.Get(url)

	if err != nil {
		return nil, err
	}

	xkcd := &XKCD{}

	if err := json.Unmarshal(result, xkcd); err != nil {
		return nil, fmt.Errorf("fetch: %v", err)
	}

	return xkcd, nil
}

// getBinarySearch is a recursive algorithm for locating the comicNum based on startIx and endIx index values.
// If the comic is found it will return the index location and an object, otherwise it returns -1 and nil.
// This method will be invoked from Get only if Comics.sorted is set to true.
func (c *Comics) getBinarySearch(startIx, endIx int, comicNum int) (int, *XKCD) {
	if startIx > endIx {
		return -1, nil
	}

	mid := (endIx + startIx) / 2

	if c.comics[mid].Number == comicNum {
		return mid, &c.comics[mid]
	} else if c.comics[mid].Number < comicNum {
		startIx = mid + 1
	} else {
		endIx = mid - 1
	}

	return c.getBinarySearch(startIx, endIx, comicNum)
}

// getSequentialSearch will scan through the Comics.comics slice from start to finish
// in order to locate the comicNum. If the comic is found it returns index location and XKCD object,
// otherwise it will return -1 and nil.
func (c *Comics) getSequentialSearch(comicNum int) (int, *XKCD) {
	for i, item := range c.comics {
		if item.Number == comicNum {
			return i, &item
		}
	}

	return -1, nil
}
