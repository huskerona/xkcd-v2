package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"xkcd2/comic"
	"xkcd2/persistence"
	"xkcd2/tools/logger"
)

var (
	logging = flag.Bool("l", false, "creates a log files")
	stat    = flag.Bool("s", false, "show offline index stats")
	dump    = flag.Bool("d", false, "output comic numbers, year, month and date (used only with -s)")
)

var (
	wg            sync.WaitGroup
	lastComicChan chan int         // last comic number found on the web site
	comicChan     chan *comic.XKCD // downloaded comic
	statusChan    <-chan time.Time // time to refresh the progress status

	comics comic.Comics
)

// func init() {
// 	comics = make(comic.Comics, 0)
// }

func main() {
	flag.Parse()
	logger.Initialize(*logging)

	defer logger.Trace("main")()

	loadComics()

	if !*stat {
		// no flag
		start := time.Now()
		doSync()
		comics.Sort()
		writeComics()

		fmt.Printf("\nDONE in %s\n", time.Since(start))
		fmt.Printf("\nTotal comics: %d\n", comics.Len())
	} else {
		// -s flag
		if *dump {
			// -d flag
			for _, item := range comics.GetAll() {
				fmt.Printf("%d,%s,%s,%s\n",
					item.Number, item.Year, item.Month, item.Day)
			}
		}

		fmt.Printf("\nIndex status: %d\n", comics.Len())
	}
}

// doSync initiates the syncing process of fetching the latest comic, fetch missing comics,
// sorting the comics and writing them back to the offline index file.
func doSync() {
	lastComicChan = make(chan int)
	comicChan = make(chan *comic.XKCD)
	statusChan = time.Tick(500 * time.Millisecond)

	go monitor()

	lastComicNum := getLatestComicNum()

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
	// to the Download function.
	semaphore := make(chan struct{}, 20)

	for i := 1; i < lastComicNum; i++ {
		if comics.Contains(i) {
			continue
		}

		wg.Add(1)

		go func(comicNum int) {
			defer logger.Trace(fmt.Sprintf("fetchComics go func(%d)", comicNum))()
			defer wg.Done()

			xkcd := &comic.XKCD{}
			var err error

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err = xkcd.Download(comicNum); err != nil {
				return
			}

			comicChan <- xkcd
		}(i)
	}
}

// Loads the comics from the index file
func loadComics() {
	defer logger.Trace("loadComics")()

	temp, err := persistence.ReadIndexFile()

	if err != nil {
		log.Println(err)
	}

	if temp != nil {
		comics.Load(temp)
	}
}

// Writes the comics back to the index file
func writeComics() {
	err := persistence.WriteIndexFile(comics.GetAll())

	if err != nil {
		log.Fatal(err)
	}
}

// Retrieves the latest comic and passes the information to lastComicChan and comicChan channels.
func getLatestComicNum() int {
	xkcd := &comic.XKCD{}
	err := xkcd.Download(0)

	if err != nil {
		log.Fatal(err)
	}

	result := xkcd.Number
	comicChan <- xkcd
	lastComicChan <- result

	return result
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
				comics.Add(item)
			}

		case <-statusChan:
			result := float64(comics.Len()) / float64(lastComicNum) * 100
			fmt.Printf("Downloading: %.2f%%\r", result)
		}
	}
}

func closeChannels() {
	close(comicChan)
}
