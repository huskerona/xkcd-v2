package webclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

// Get makes a call to the web server and returns a byte slice with raw data.
// The calling function should handle the slice either by decoding/unmarshalling the JSON or
// doing something else with the byte slice.
func Get(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("NewRequest: %v", err)

	}
	resp, err := Client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Get: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status: %v", resp.Status)
	}

	var result []byte

	if result, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("ReadAll: %v", err)
	}

	return result, nil
}
