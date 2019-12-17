package util

import (
	"log"
	"os/user"

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
