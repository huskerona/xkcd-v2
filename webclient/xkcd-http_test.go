package webclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"xkcd2/webclient/mocks"
)

func TestGetOK(t *testing.T) {
	want := "This is a test"
	Client = &mocks.MockClient{}
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		data := ioutil.NopCloser(bytes.NewReader([]byte(want)))

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       data,
		}, nil
	}

	res, err := Get("test-url")

	if err != nil {
		t.Errorf("expected err to be nil, got %v", err)
	}

	got := string(res)

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}

}

func TestGetBadRequest(t *testing.T) {
	Client = &mocks.MockClient{}
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		data := ioutil.NopCloser(bytes.NewReader([]byte("")))

		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       data,
		}, nil
	}

	_, err := Get("test-url")

	if err == nil {
		t.Errorf("expected invalid response status")
	}
}
