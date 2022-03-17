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
}

// Load will add comics to c. During the loading process the critical section is formed before looping over the comics.
func (c *Comics) Load(items []XKCD) {
	defer logger.Trace("func LoadComics")()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Len() == 0 {
		for _, item := range items {
			c.comics = append(c.comics, item)
		}
	}
}

// Add will insert xkcd into a collection of comics. Add uses a mutex to add an item
// in order to prevent concurrent access to the collection.
func (c *Comics) Add(xkcd *XKCD) {
	defer logger.Trace("method Add()")()

	c.mu.Lock()
	defer c.mu.Unlock()
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

	index := -1
	var xkcd *XKCD

	c.mu.Lock()
	defer c.mu.Unlock()

	for i, item := range c.comics {
		if item.Number == comicNum {
			index, xkcd = i, &item
			break
		}
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
func (c *Comics) Remove(index int) {
	defer logger.Trace("method Remove()")()

	if index == -1 || index > len(c.comics) {
		return
	}

	copy(c.comics[index:], c.comics[index+1:])
}

// Sort will sort the comics in ascending order based on the comic number
func (c *Comics) Sort() {
	c.mu.Lock()
	defer c.mu.Unlock()
	sort.Sort(c)
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
