package main

import (
	"flag"
	"fmt"

	// "github.com/huskerona/xkcd2/infrastructure/util"
	appStr "github.com/huskerona/xkcd2/infrastructure/app-strings"
	"github.com/huskerona/xkcd2/infrastructure/logger"
)

var log = flag.Bool("log", false, "creates a log files")
var verbose = flag.Bool("v", false, "verbose output")

func main() {
	flag.Parse()

	fmt.Println(appStr.AppTitle)

	logger.Initialize(*log)

	defer logger.Trace("main")()

	fmt.Println("exit...")
}