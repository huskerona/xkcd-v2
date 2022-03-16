package comic

import (
	"fmt"
	"xkcd2/config"
	"xkcd2/tools/imaging"
	"xkcd2/tools/logger"
)

//+ Type defs
type XKCD struct {
	Day        string `json:"day"`
	Month      string `json:"month"`
	Year       string `json:"year"`
	Number     int    `json:"num"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	ImageURL   string `json:"img"`
	ImageAlt   string `json:"alt"`
	News       string `json:"news"`
	Link       string `json:"link"`

	// base64 encoded image downloaded from ImageURL
	Image string
}

// Download fetches the JSON contents of the XKCD comic based on its number. If number is 0 it will
// fetch the latest issue. The JSON file is unmarshalled into XKCD structure. The Image file is not downloaded
// and it needs a separate call to DownloadImage
func (xkcd *XKCD) Download(comicNumber int) error {
	defer logger.Trace("func DownloadComic")()

	var err error
	var url string

	if comicNumber == 0 {
		url = fmt.Sprintf("%s/%s", config.HomeURL, config.JSONURL)
	} else {
		url = fmt.Sprintf("%s/%d/%s", config.HomeURL, comicNumber, config.JSONURL)
	}

	xkcd, err = fetch(url)

	if err != nil {
		return err
	}

	return nil
}

// DownloadImage fetches an XKCD image from imageUrl.
// NOTE: There are some comics whose image cannot be retrieved. It would require that we parse the HTML.
// For that reason the error is ignored, but in any case, XKCD struct is returned while Image is left as empty string.
func (xkcd *XKCD) DownloadImage(imageUrl string) error {
	imageByte, err := downloadImage(imageUrl)

	if err == nil {
		xkcd.Image = imaging.EncodeToBase64(imageByte)
	}

	return nil
}
