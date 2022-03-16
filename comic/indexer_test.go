package comic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"xkcd2/webclient"
	"xkcd2/webclient/mocks"
)

func setupClient(url string, forceError bool) *mocks.MockClient {
	want := fmt.Sprintf(`{
		"day": "1", 
		"month": "1", 
		"year": "2022", 
		"num": 1, 
		"title": "unit test", 
		"safe_title": "safe unit test", 
		"transcript": "", 
		"img": "%s",
		"alt": "unit test",
		"news": "",
		"link": "http://localhost/1"
	}`, url)

	mockClient := &mocks.MockClient{}
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		data := ioutil.NopCloser(bytes.NewReader([]byte(want)))

		var err error

		if forceError {
			err = fmt.Errorf("unit test error")
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       data,
		}, err
	}

	return mockClient
}

func TestFetch(t *testing.T) {
	url := "http://localhost/1/image.jpg"
	webclient.Client = setupClient(url, false)

	got, err := fetch("test-url")

	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}

	if got.Title != "unit test" {
		t.Errorf("expected title %s, got %s", "unit test", got.Title)
	}

	if got.ImageURL != url {
		t.Errorf("expected image url %s, got %s", url, got.ImageURL)
	}
}

func TestFetchError(t *testing.T) {
	url := "http://localhost/1/image.jpg"
	webclient.Client = setupClient(url, true)

	_, err := fetch("test-url")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDownloadComicLatest(t *testing.T) {
	xkcd := XKCD{}

	url := "http://localhost/1/image.jpg"
	webclient.Client = setupClient(url, false)

	err := xkcd.Download(0)

	if err != nil {
		t.Errorf("expected err to be nil, got %v", err)
	}

	if xkcd.Image == "" {
		t.Error("expected xkcd.Image to contain value, got empty string")
	}
}

func TestDownloadComicOlder(t *testing.T) {

}
