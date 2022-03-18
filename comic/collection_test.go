package comic

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func setupComics(total int, randomize bool) []XKCD {
	comics := make([]XKCD, 0, total)
	var rng *rand.Rand

	if randomize {
		seed := time.Now().UnixMilli()
		rng = rand.New(rand.NewSource(seed))
	}

	for i := 1; i <= total; i++ {
		id := 0

		if randomize {
			id = rng.Intn(total + 1)
		} else {
			id = i
		}

		xkcd := XKCD{Number: id, Title: fmt.Sprintf("Comic %d", i), Link: fmt.Sprintf("https://xkcd.com/%d/info.0.json", i)}
		comics = append(comics, xkcd)
	}

	return comics
}

func TestComicsLoad(t *testing.T) {
	want := 100
	comics := setupComics(want, false)

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

func TestComicsAddUnordered(t *testing.T) {
	c := Comics{}

	var xkcd10 = &XKCD{Number: 10, Title: "Adding Title"}
	var xkcd9 = &XKCD{Number: 9, Title: "Adding Title"}

	c.sorted = true

	c.Add(xkcd10)
	c.Add(xkcd9)

	got := c.sorted

	if got {
		t.Errorf("wanted sorted=false, got %t", got)
	}
}

func TestComicsGet(t *testing.T) {
	c := Comics{}
	c.Load(setupComics(2000, false))

	comicNum := 700
	index := comicNum - 1

	i, xkcd := c.Get(comicNum)

	if xkcd == nil {
		t.Error("expected object, got nil")
	}

	if i != index {
		t.Errorf("expected %d, got %d", index, i)
	}
}

func TestComicsGetFirst(t *testing.T) {
	c := Comics{}
	c.Load(setupComics(100, false))

	comicNum := 1
	index := comicNum - 1

	i, xkcd := c.Get(comicNum)

	if xkcd == nil {
		t.Error("expected object, got nil")
	}

	if i != index {
		t.Errorf("expected %d, got %d", index, i)
	}
}

func TestComicsGetLast(t *testing.T) {
	c := Comics{}
	c.Load(setupComics(100, false))

	comicNum := 100
	index := comicNum - 1

	i, xkcd := c.Get(comicNum)

	if xkcd == nil {
		t.Error("expected object, got nil")
	}

	if i != index {
		t.Errorf("expected %d, got %d", index, i)
	}
}

func TestComicsGetByComicNum100000(t *testing.T) {
	c := Comics{}
	c.Load(setupComics(100, false))

	comicNum := 100000

	i, xkcd := c.Get(comicNum)

	if i != -1 {
		t.Errorf("expected -1, got %d", i)
	}

	if xkcd != nil {
		t.Errorf("expected nil, got %v", xkcd)
	}
}

func TestComicsGetByComicNum0(t *testing.T) {
	c := Comics{}
	c.Load(setupComics(100, false))

	comicNum := 0

	i, xkcd := c.Get(comicNum)

	if i != -1 {
		t.Errorf("expected -1, got %d", i)
	}

	if xkcd != nil {
		t.Errorf("expected nil, got %v", xkcd)
	}
}

func TestRemoveExisting(t *testing.T) {
	c := Comics{}
	c.Load(setupComics(100, false))

	index := 50

	want := len(c.comics)
	ok := c.Remove(index)
	got := len(c.comics)

	if !ok {
		t.Errorf("expected true, got %t", ok)
	}

	if want == got {
		t.Errorf("expected %d, got %d", want, got)
	}
}

func TestRemoveMinus1(t *testing.T) {
	c := Comics{}
	c.Load(setupComics(10, false))

	index := -1

	ok := c.Remove(index)

	if ok {
		t.Errorf("want false, got %t", ok)
	}
}

func TestRemoveOverLength(t *testing.T) {
	c := Comics{}
	c.Load(setupComics(10, false))

	index := len(c.comics)

	ok := c.Remove(index)

	if ok {
		t.Errorf("want false, got %t", ok)
	}
}

func TestSort(t *testing.T) {
	c := Comics{}
	c.Load(setupComics(1000, true))
	c.sorted = false // input is randomized so we need to change to false

	want := c.sorted
	c.Sort()
	got := c.sorted

	if want == got {
		t.Errorf("expected %t, got %t", want, got)
	}

	for i := 1; i < c.Len(); i++ {
		prev := c.comics[i-1].Number
		curr := c.comics[i].Number

		if prev > curr {
			t.Errorf("expected %d, got %d", curr, prev)
		}
	}
}
