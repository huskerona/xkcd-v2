/*
Package comic is used for working with single or a collection of XKCD comics.

Basics

When working with a single comic, the XKCD type used the following methods to update the information:

    Download(comicNumber int) error
    DownloadImage(imageUrl string) error

Download method will determine which comic to fetch based on the comicNumber. If the comicNumber is 0,
the latest version of the comic will be fetched using https://xkcd.com/info.0.json, however, if the
comicNumber is greater than 0, the URL will be, for example, https://xkcd.com/123/info.0.json.
When the JSON is retrieved, it will contain the information about the image URL. The image is not downloaded
during this process.

DownloadImage method uses the image URL from the img element in the JSON file (or XKCD.ImageURL) to download
the image and convert it into a Base-64 encoded string. This string is stored in XKCD.Image. Some comics are
returned with an img but the image cannot be retrieved. It requires parsing the HTML page to find the exact
URL of the image. For that reason DownloadImage will always return a nil value for error.

Types and Values

The key type in comic package is XKCD struct that stores unmarshalled data from the JSON document retrieved
from the xkcd website.

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

The package also defines a collection with which to work.

    type Comics []XKCD

The type Comics has a number of methods that allows it to manipulate the collection of XKCD values.

Working with Comics collection

The Comics collection can load, retrieve, remove and sort itself. When working with the collection
for the first time, you can either use Load(comics []XKCD) or Add(xkcd *XKCD) methods. The first method
is used when you have comics stored in the file and you want to load them into a collection. For more
information on loading see package persistence. The other method is used when you need to add one by one
comic to the collection, as in the case when a new comic is added after XKCD.Download and XKCD.DownloadImage.


*/
package comic
