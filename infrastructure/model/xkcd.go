package model

import (
	"github.com/huskerona/xkcd2/infrastructure/logger"
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

//- Type defs

//+ Exports
var Comics collection

//- Exports

//+ Methods

// Add new XKCD comic to the collection
func (c collection) Add(xkcd *XKCD) {
	defer logger.Trace("method Add()")()

	if !c.Contains((*xkcd).Number) {
		Comics = append(Comics, xkcd)
	}
}

// Determines whether a given comic is in the collection
func (c collection) Contains(comicNum int) bool {
	defer logger.Trace("method Contains()")()

	xkcd := c.Get(comicNum)

	return xkcd != nil
}

// Retrieves the comic from the collection. If there is no such comic, nil is returned.
func (c collection) Get(comicNum int) *XKCD {
	defer logger.Trace("method Get()")()

	for _, item := range c {
		if item.Number == comicNum {
			return item
		}
	}

	return nil
}

// Finds an index of the comic based on the comic number. If nothing was found, -1 is returned.
func (c collection) Index(comicNum int) int {
	result := -1

	for index, item := range c {
		if item.Number == comicNum {
			result = index
			break
		}
	}

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
