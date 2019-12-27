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

var (
	logging = flag.Bool("log", false, "creates a log files")
	stat    = flag.Bool("stat", false, "show offline index stats")
	dump    = flag.Bool("dump", false, "output comic numbers, year, month and date (only with -stat)")
)

var (
	wg                 sync.WaitGroup
	done               chan bool        // signals the end of the operations
	lastComicChan      chan int         // last comic number found on the web site
	comicChan          chan *model.XKCD // downloaded comic
	statusChan         <-chan time.Time // time to refresh the progress status
	duplicateCheckChan chan duplicate   // duplicate comic check channel
)

// contains the information about the duplicate request.
// comicNum field carries the number of the comic to check, and
// dupCheckChan is the channel on which the monitor goroutine will
// return the information if the comic is duplicate or not.
type duplicate struct {
	comicNum     int
	dupCheckChan chan bool
}

// initialization of the channels
func init() {
	lastComicChan = make(chan int)
	comicChan = make(chan *model.XKCD)
	duplicateCheckChan = make(chan duplicate)
	statusChan = time.Tick(500 * time.Millisecond)
	done = make(chan bool)
}

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
				_, _ = fmt.Printf("%d,%s,%s,%s\n",
					(*item).Number, (*item).Year, (*item).Month, (*item).Day)
			}
		}

		fmt.Printf("\nIndex status: %d\n", indexManager.Count())
	}
}

// Initiates the syncing process of fetching the latest comic, fetch missing comics,
// sorting the comics and writing them back to the offline index file.
func doSync() {
	go monitor()

	lastComicNum := getLastComicNum()

	fetchComics(lastComicNum)

	if !completed() {
		time.Sleep(1 * time.Second)
	}

	indexManager.Sort()
	writeComics()
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

	// Channel closer
	go func() {
		wg.Wait()
		done <- true
	}()
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
	result := make(chan bool)

	dc := duplicate{comicNum: comicNum, dupCheckChan: result}
	duplicateCheckChan <- dc

	return <-result
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

		case dupCheck := <-duplicateCheckChan:
			dupCheck.dupCheckChan <- indexManager.Contains(dupCheck.comicNum)

		case <-statusChan:
			result := float64(indexManager.Count()) / float64(lastComicNum) * 100
			fmt.Printf("Downloading: %.2f%%\r", result)

		case complete := <-done:
			if complete {
				fmt.Println("Exiting...")
				closeChannels()
				return
			}
		}
	}
}

// Checks if the operation is completed and re-sends the received value to the channel done.
func completed() bool {
	isDone := <-done
	done <- isDone

	return isDone
}

func closeChannels() {
	close(comicChan)
	close(duplicateCheckChan)
	close(done)
}