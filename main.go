package main

import (
	"flag"
	"fmt"
	"github.com/huskerona/xkcd2/infrastructure/model"
	"log"
	"os"
	"sync"
	"time"

	fileManager "github.com/huskerona/xkcd2/file-manager"
	indexManager "github.com/huskerona/xkcd2/index-manager"
	"github.com/huskerona/xkcd2/infrastructure/logger"
)

var logging = flag.Bool("log", false, "creates a log files")
var verbose = flag.Bool("v", false, "verbose output")
var stat = flag.Bool("stat", false, "show offline index stats")
var dump = flag.Bool("dump", false, "output comic numbers, year, month and date (only with -stat)")

var wg sync.WaitGroup

func main() {
	flag.Parse()
	logger.Initialize(*logging)

	defer logger.Trace("main")()

	loadComics()

	if !*stat {
		start := time.Now()
		doSync()
		fmt.Printf("\nDONE in %s\n", time.Since(start))

		fmt.Printf("\nTotal comics: %d\n", indexManager.Count())
	} else {
		if *dump {
			for _, item := range indexManager.GetComics() {
				fmt.Fprintf(os.Stdout, "%d,%s,%s,%s\n",
					(*item).Number, (*item).Year, (*item).Month, (*item).Day)
				//fmt.Fprintf(os.Stdout, "%v\n", item)
			}
		}

		fmt.Printf("\nIndex status: %d\n", indexManager.Count())
	}
}

func doSync() {
	//defer wg.Done()
	lastComicNum := getLastComicNum()

	go spinner()
	synchronize(lastComicNum)

	indexManager.Sort()
	writeComics()
}

func synchronize(lastComicNum int) {
	defer logger.Trace("synchronize")()

	ch := make(chan *model.XKCD)
	// counting semaphore token that enforces the limit on the number of calls
	// to the DownloadComic function.
	semaphore := make(chan struct{}, 20)

	n := 0 // Terminates the loop which adds xkcd to the index file.

	for i := 1; i < lastComicNum; i++ {
		if indexManager.Contains(i) {
			continue
		}

		wg.Add(1)
		n++

		go func(comicNum int) {
			defer logger.Trace(fmt.Sprintf("synchronize go func(%d)", comicNum))()
			defer wg.Done()

			var xkcd *model.XKCD
			var err error
			hasError := false

			semaphore <- struct{}{}
			if xkcd, err = indexManager.DownloadComic(comicNum); err != nil {
				hasError = true
			}

			if !hasError {

				ch <- xkcd
			}
			<-semaphore
			n--
		}(i)
	}

	// Channel closer
	go func() {
		wg.Wait()

		close(ch)
	}()

	// Previous implementation used for item := range ch to empty the channel and to
	// add comics to the collection but the problem was that on some occasions the
	// channel was not emptied fully. This sort of approach is described in
	// The Go Programming Language, p.242
	for n > 0 {
		indexManager.AddToCollection(<-ch)
	}
}

func spinner() {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("Downloading... %c\r", r)
			time.Sleep(100 * time.Millisecond)
		}
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
