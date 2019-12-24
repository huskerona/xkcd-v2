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
var stat = flag.Bool("stat", false, "show offline index stats")
var dump = flag.Bool("dump", false, "output comic numbers, year, month and date (only with -stat)")
var showDownloads = flag.Bool("sd", false, "show status of each comic as it is being downloaded")

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
				_, _  = fmt.Printf("%d,%s,%s,%s\n",
					(*item).Number, (*item).Year, (*item).Month, (*item).Day)
				//fmt.Fprintf(os.Stdout, "%v\n", item)
			}
		}

		fmt.Printf("\nIndex status: %d\n", indexManager.Count())
	}
}

func doSync() {
	lastComicNum := getLastComicNum()

	if !*showDownloads {
		go progressCounter()
	}

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

	for i := 1; i < lastComicNum; i++ {
		if indexManager.Contains(i) {
			continue
		}

		wg.Add(1)

		go func(comicNum int) {
			defer logger.Trace(fmt.Sprintf("synchronize go func(%d)", comicNum))()
			defer wg.Done()

			var xkcd *model.XKCD
			var err error
			hasError := false

			semaphore <- struct{}{}
			if xkcd, err = indexManager.DownloadComic(comicNum); err != nil {
				hasError = true
				// Without hasError variable here, any error that occurs would
				// not remove the token from the semaphore if we just used return
			}

			if !hasError {
				if *showDownloads {
					fmt.Printf("Comic %d of %d\n", comicNum, lastComicNum)
				}

				ch <- xkcd
			}
			<-semaphore // release the token
		}(i)
	}

	// Channel closer
	go func() {
		wg.Wait()
		close(ch)
	}()

	for item := range ch {
		indexManager.AddToCollection(item)
	}
}

func progressCounter() {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("Downloading...%c\r", r)
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
