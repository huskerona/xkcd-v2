package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	fileManager "xkcd2/file-manager"
	indexManager "xkcd2/index-manager"
	"xkcd2/infrastructure/logger"
	"xkcd2/infrastructure/model"
)

var (
	logging = flag.Bool("l", false, "creates a log files")
	stat    = flag.Bool("s", false, "show offline index stats")
	dump    = flag.Bool("d", false, "output comic numbers, year, month and date (used only with -s)")
)

var (
	wg            sync.WaitGroup
	lastComicChan chan int         // last comic number found on the web site
	comicChan     chan *model.XKCD // downloaded comic
	statusChan    <-chan time.Time // time to refresh the progress status
)

func main() {
	flag.Parse()
	logger.Initialize(*logging)

	defer logger.Trace("main")()

	loadComics()

	if !*stat {
		// no flag
		start := time.Now()
		doSync()
		indexManager.Sort()
		writeComics()

		fmt.Printf("\nDONE in %s\n", time.Since(start))
		fmt.Printf("\nTotal comics: %d\n", indexManager.Count())
	} else {
		// -s flag
		if *dump {
			// -d flag
			for _, item := range indexManager.GetComics() {
				fmt.Printf("%d,%s,%s,%s\n",
					(*item).Number, (*item).Year, (*item).Month, (*item).Day)
			}
		}

		fmt.Printf("\nIndex status: %d\n", indexManager.Count())
	}
}

// Initiates the syncing process of fetching the latest comic, fetch missing comics,
// sorting the comics and writing them back to the offline index file.
func doSync() {
	lastComicChan = make(chan int)
	comicChan = make(chan *model.XKCD)
	statusChan = time.Tick(500 * time.Millisecond)

	go monitor()

	lastComicNum := getLastComicNum()

	fetchComics(lastComicNum)

	// Channel closer
	wg.Wait()
	closeChannels()
}

// fetchComics function does the actual hard work of downloading all the missing comics.
// The lastComicNum parameter is the latest comic on the XKCD web site.
func fetchComics(lastComicNum int) {
	defer logger.Trace("fetchComics")()

	// counting semaphore token that enforces the limit on the number of calls
	// to the DownloadComic function.
	semaphore := make(chan struct{}, 20)

	for i := 1; i < lastComicNum; i++ {
		if comicExists(i) {
			continue
		}

		wg.Add(1)

		go func(comicNum int) {
			defer logger.Trace(fmt.Sprintf("fetchComics go func(%d)", comicNum))()
			defer wg.Done()

			var xkcd *model.XKCD
			var err error

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if xkcd, err = indexManager.DownloadComic(comicNum); err != nil {
				return
			}

			comicChan <- xkcd
		}(i)
	}
}

// Loads the comics from the index file
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

// Writes the comics back to the index file
func writeComics() {
	err := fileManager.WriteIndexFile(indexManager.GetComics())

	if err != nil {
		log.Fatal(err)
	}
}

// Retrieves the latest comic and passes the information to lastComicChan and comicChan channels.
func getLastComicNum() int {
	xkcd, err := indexManager.DownloadComic(0)

	if err != nil {
		log.Fatal(err)
	}

	result := (*xkcd).Number
	comicChan <- xkcd
	lastComicChan <- result

	return result
}

// Checks if the comic exists by sending the information to
// the monitor goroutine and receives the feedback on the dupCheckChan channel.
func comicExists(comicNum int) bool {
	return indexManager.Contains(comicNum)
}

// monitor function monitors the channels and does something with the
// data that arrives on each channel.
func monitor() {
	var lastComicNum int

	for {
		select {
		case comicNum, ok := <-lastComicChan:
			if ok {
				lastComicNum = comicNum
				close(lastComicChan)
			}

		case item := <-comicChan:
			if item != nil {
				logger.Info(fmt.Sprintf("Adding %d\n", (*item).Number))
				indexManager.AddToCollection(item)
			}

		case <-statusChan:
			result := float64(indexManager.Count()) / float64(lastComicNum) * 100
			fmt.Printf("Downloading: %.2f%%\r", result)
		}
	}
}

func closeChannels() {
	close(comicChan)
}
