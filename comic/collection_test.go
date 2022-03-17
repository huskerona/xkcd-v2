package comic

import (
	"fmt"
	"testing"
)

func setupComics(total int) []XKCD {
	comics := make([]XKCD, 0, total)
	for i := 1; i <= total; i++ {
		xkcd := XKCD{Number: i, Title: fmt.Sprintf("Comic %d", i), Link: fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)}
		comics = append(comics, xkcd)
	}

	return comics
}

func TestComicsLoad(t *testing.T) {
	want := 100
	comics := setupComics(want)

	c := Comics{}
	c.Load(comics)

	if c.Len() != len(comics) {
		t.Errorf("wanted %d, got %d", len(comics), c.Len())
	}
}

func TestComicsAdd(t *testing.T) {
	c := Comics{}

	var xkcd = &XKCD{Number: 1, Title: "Adding Title"}

	c.Add(xkcd)

	if c.Len() != 1 {
		t.Errorf("wanted 1, got %d", c.Len())
	}
}
