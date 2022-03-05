package util

import (
	"fmt"
	"log"
	"os/user"

	appStrings "xkcd2/infrastructure/app-strings"
)

// Defines the home folder where XKCD will be written to.
func GetHomeFolder() string {
	var result string

	user, err := user.Current()

	if err != nil {
		log.Fatalf("gethomefolder: %v", err)
	}

	result = user.HomeDir

	return result
}

// Returns the location of the XKCD folder
func GetXkcdFolder() string {
	return fmt.Sprintf("%s/.xkcd", GetHomeFolder())
}

// Returns complete filename of the XKCD index file
func GetIndexFile() string {
	return fmt.Sprintf("%s/%s", GetXkcdFolder(), appStrings.IndexFile)
}
