package main

import (
	"flag"
	"fmt"
	"github.com/huskerona/xkcd2/infrastructure/model"
	"log"
	"sync"
	"time"

	fileManager "github.com/huskerona/xkcd2/file-manager"
	indexManager "github.com/huskerona/xkcd2/index-manager"
	"github.com/huskerona/xkcd2/infrastructure/logger"
)

var logging = flag.Bool("log", false, "creates a log files")
var verbose = flag.Bool("v", false, "verbose output")
var stat = flag.Bool("stat", false, "show offline index stats")

var wg sync.WaitGroup

func main() {
	flag.Parse()
	logger.Initialize(*logging)

	defer logger.Trace("main")()

	start := time.Now()

	//wg.Add(1)
	doSync()

	//wg.Wait()

	fmt.Printf("DONE in %s\n", time.Since(start))
	fmt.Printf("\nIndex status: %d\n", indexManager.Count())
}

func doSync() {
	//defer wg.Done()
	wg.Add(1)
	go func() {
		defer wg.Done()
		go loadComics()
	}()

	wg.Wait()

	lastComicNum := getLastComicNum()

	synchronize(lastComicNum)

	indexManager.Sort()
	writeComics()
}

func synchronize(lastComicNum int) {
	defer logger.Trace("synchronize")()

	ch := make(chan *model.XKCD)

	for i := 1; i < lastComicNum; i++ {
		if indexManager.Contains(i) {
			continue
		}

		wg.Add(1)

		go func(i int) {
			defer logger.Trace(fmt.Sprintf("synchronize go func(%d)", i))()
			defer wg.Done()

			var xkcd *model.XKCD
			var err error

			if xkcd, err = indexManager.DownloadComic(i); err != nil {
				return
			}

			ch <- xkcd

			fmt.Printf("Downloading %d of %d\n", i, lastComicNum)
		}(i)
	}

	// Channel closer
	go func() {
		wg.Wait()
		close(ch)
		fmt.Println("Channel closed")
	}()

	for item := range ch {
		indexManager.AddToCollection(item)
	}
}

func loadComics() {
	defer logger.Trace("loadComics")()

	comics, err := fileManager.ReadIndexFile()

	if err != nil {
		log.Println(err)
	}

	if comics != nil {
		indexManager.LoadComics(comics)
	}
}

func writeComics() {
	err := fileManager.WriteIndexFile(indexManager.GetComics())

	if err != nil {
		log.Fatal(err)
	}
}

func getLastComicNum() int {
	xkcd, err := indexManager.DownloadComic(0)

	if err != nil {
		log.Fatal(err)
	}

	indexManager.AddToCollection(xkcd)

	result := (*xkcd).Number

	return result
}
