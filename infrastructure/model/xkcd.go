package model

import (
	//"sync"

	"github.com/huskerona/xkcd2/infrastructure/logger"
	"github.com/sasha-s/go-deadlock"
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

type collection []*XKCD
var mu deadlock.Mutex

//- Type defs

//+ Exports
var Comics collection

//- Exports

//+ Methods

// Add new XKCD comic to the collection
func (c collection) Add(xkcd *XKCD) {
	defer logger.Trace("method Add()")()

	if !c.Contains((*xkcd).Number) {
		mu.Lock()
		Comics = append(Comics, xkcd)
		mu.Unlock()
	}
}

// Determines whether a given comic is in the collection
func (c collection) Contains(comicNum int) bool {
	defer logger.Trace("method Contains()")()

	index, _ := c.Get(comicNum)

	return index > -1
}

// Retrieves the comic from the collection together with the index where the comic was found.
//If there is no such comic, -1, nil are returned.
func (c collection) Get(comicNum int) (int, *XKCD) {
	defer logger.Trace("method Get()")()

	index := -1
	var xkcd *XKCD

	mu.Lock()
	defer mu.Unlock()

	for i, item := range c {
		if item.Number == comicNum {
			index, xkcd = i, item
			break
		}
	}

	return index, xkcd
}

// Finds an index of the comic based on the comic number. If nothing was found, -1 is returned.
func (c collection) Index(comicNum int) int {
	result, _ := c.Get(comicNum)

	return result
}

// Removes an item from the collection based on its index
func (c collection) Remove(index int) {
	defer logger.Trace("method Remove()")()

	if index == -1 || index > len(Comics) {
		return
	}

	copy(Comics[index:], Comics[index+1:])
}

//++ Methods (sort interface)
// Returns the length of the collection
func (c collection) Len() int {
	return len(c)
}

// Returns true if item at index i is small than item at index j
func (c collection) Less(i, j int) bool {
	return c[i].Number < c[j].Number
}

// Swaps two items in the collection
func (c collection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

//-- Methods (sort interface)
//- Methods
