package comic

import (
	"encoding/json"
	"fmt"
	"xkcd2/tools/logger"
	"xkcd2/webclient"
)

// downloadImage fetches an image specified in imageUrl parameter.
func downloadImage(imageUrl string) ([]byte, error) {
	defer logger.Trace("func downloadImage()")()

	result, err := webclient.Get(imageUrl)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// fetch perfroms GET operation on url and it unmarshalls the JSON document into XKCD object.
func fetch(url string) (*XKCD, error) {
	result, err := webclient.Get(url)

	if err != nil {
		return nil, err
	}

	var xkcd *XKCD
	if err := json.Unmarshal(result, &xkcd); err != nil {
		return nil, fmt.Errorf("fetch: %v", err)
	}

	return xkcd, nil
}
