package comic

import (
	"fmt"
	"sort"
	"sync"
	"xkcd2/tools/logger"
)

type Comics struct {
	mu     sync.Mutex
	comics []XKCD
	sorted bool
}

// Load will add comics to c. During the loading process the critical section is formed before looping over the comics.
// When the data is loaded, an assumption is made that the items are sorted, therefore, prior to writing the comics
// to the disk, make sure to call Sort() method.
func (c *Comics) Load(items []XKCD) {
	defer logger.Trace("func LoadComics")()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Len() == 0 {
		c.comics = make([]XKCD, 0, len(items))
		c.comics = append(c.comics, items...)
		c.sorted = true
	}
}

// Add will insert xkcd into a collection of comics. Add uses a mutex to add an item
// in order to prevent concurrent access to the collection.
func (c *Comics) Add(xkcd *XKCD) {
	defer logger.Trace("method Add()")()

	c.mu.Lock()
	defer c.mu.Unlock()

	size := c.Len()
	if size > 0 && c.comics[size-1].Number > xkcd.Number {
		c.sorted = false
	}

	c.comics = append(c.comics, *xkcd)
}

// Contains return true or false depending on whether it found a comicNum in the collection
func (c *Comics) Contains(comicNum int) bool {
	defer logger.Trace("method Contains()")()

	index, _ := c.Get(comicNum)

	logger.Info(fmt.Sprintf("%d exists: %t", comicNum, index > -1))
	return index > -1
}

// Get returns a comic and the index where the comicNum was found. If a comic is not found, the return value is -1 for index and nil for XKCD type
func (c *Comics) Get(comicNum int) (int, *XKCD) {
	defer logger.Trace("method Get()")()

	c.mu.Lock()
	defer c.mu.Unlock()

	index := -1
	var xkcd *XKCD

	if c.sorted {
		index, xkcd = c.getBinarySearch(0, len(c.comics)-1, comicNum)
	} else {
		index, xkcd = c.getSequentialSearch(comicNum)
	}

	return index, xkcd
}

// GetAll returns the entire collection of Comics
func (c *Comics) GetAll() []XKCD {
	defer logger.Trace("func GetComics")()

	return c.comics
}

// Index returns index number of a collection or -1 if not found. See Get func.
func (c *Comics) Index(comicNum int) int {
	result, _ := c.Get(comicNum)

	return result
}

// Remove drops an item from the collection based on its index
func (c *Comics) Remove(index int) bool {
	defer logger.Trace("method Remove()")()

	if index == -1 || index >= len(c.comics) {
		return false
	}

	copy(c.comics[index:], c.comics[index+1:])
	c.comics = c.comics[:len(c.comics)-1]

	return true
}

// Sort will sort the comics in ascending order based on the comic number
func (c *Comics) Sort() {
	c.mu.Lock()
	defer c.mu.Unlock()
	sort.Sort(c)
	c.sorted = true
}

// Len returns a len value of the collection slice. It is part of the sort packages Interface type.
func (c *Comics) Len() int {
	return len(c.comics)
}

// Less compares i and j and returns bool if i is smaller than j.
func (c *Comics) Less(i, j int) bool {
	return c.comics[i].Number < c.comics[j].Number
}

// Swap replaces the position of two items in the collection
func (c *Comics) Swap(i, j int) {
	c.comics[i], c.comics[j] = c.comics[j], c.comics[i]
}
