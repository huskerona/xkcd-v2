package net_manager

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Makes a call to the web server and returns a byte slice with raw data.
// The calling function should handle the slice either by decoding/unmarshalling the JSON or
// doing something else with the byte slice.
func GetDataFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("getDataFromURL: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status error: %v", resp.Status)
	}

	var result []byte

	if result, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("getDataFromURL: %v", err)
	}

	return result, nil
}
