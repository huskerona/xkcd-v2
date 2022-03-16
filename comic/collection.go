package comic

import (
	"fmt"
	"sync"
	"xkcd2/tools/logger"
)

type collection []*XKCD

var (
	mu     sync.Mutex
	Comics collection
)

// Add will insert xkcd into a collection of comics. Add uses a mutex to add an item
// in order to prevent concurrent access to the collection.
func (c collection) Add(xkcd *XKCD) {
	defer logger.Trace("method Add()")()

	mu.Lock()
	defer mu.Unlock()
	Comics = append(Comics, xkcd)
}

// Contains return true or false depending on whether it found a comicNum in the collection
func (c collection) Contains(comicNum int) bool {
	defer logger.Trace("method Contains()")()

	index, _ := c.Get(comicNum)

	logger.Info(fmt.Sprintf("%d exists: %t", comicNum, index > -1))
	return index > -1
}

// Get returns a comic and the index where the comicNum was found. If a comic is not found, the return value is -1 for index and nil for XKCD type
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

// Index returns index number of a collection or -1 if not found. See Get func.
func (c collection) Index(comicNum int) int {
	result, _ := c.Get(comicNum)

	return result
}

// Remove drops an item from the collection based on its index
func (c collection) Remove(index int) {
	defer logger.Trace("method Remove()")()

	if index == -1 || index > len(Comics) {
		return
	}

	copy(c[index:], c[index+1:])
}

// Len returns a len value of the collection slice. It is part of the sort packages Interface type.
func (c collection) Len() int {
	return len(c)
}

// Less compares i and j and returns bool if i is smaller than j.
func (c collection) Less(i, j int) bool {
	return c[i].Number < c[j].Number
}

// Swap replaces the position of two items in the collection
func (c collection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
