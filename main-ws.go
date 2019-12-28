package main

import (
	"fmt"
	"log"
	"net/http"

	indexManager "github.com/huskerona/xkcd2/index-manager"
	"github.com/huskerona/xkcd2/shared"
)

func main() {
	go loadComics()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to XKCD-WS")
	})
	http.HandleFunc("/api/latest", getTotalComics)
	http.HandleFunc("/api/comic", getComic)

	if err := http.ListenAndServe("localhost:4000", nil); err != nil {
		log.Fatalf("xkcd web server: %v\n", err)
	}
}

// Loads the comics into the index manager.
func loadComics() {
	comics, err := shared.LoadComics()

	if err != nil {
		log.Fatalf("xkcd web server: %v\n", err)
	}

	indexManager.LoadComics(comics)
}

func getTotalComics(w http.ResponseWriter, req *http.Request) {
	comics := indexManager.GetComics()
	length := len(comics)
	latestComic := comics[length-1]

	if latestComic == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Errorf("xkcd-ws: comics not found\n")
		return
	}

	fmt.Fprintf(w, "Latest comic is %d", latestComic.Number)
}

func getComic(w http.ResponseWriter, req *http.Request) {

}