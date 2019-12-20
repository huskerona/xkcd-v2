package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	appStr "github.com/huskerona/xkcd2/infrastructure/app-strings"
	"github.com/huskerona/xkcd2/infrastructure/util"
)

var logWriter io.Writer
var logFile string

func init() {
	logWriter = ioutil.Discard
	logFile = fmt.Sprintf("%s/.xkcd/%s", util.GetHomeFolder(), appStr.LogFileName)
}

// Initializes the logger by defining the output file if useLog is true
func Initialize(useLog bool) {
	if !useLog {
		return
	}

	var tempWriter io.Writer

	tempWriter, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)

	if err != nil {
		log.Printf("trace initialize: %v", err)
	}

	if f, ok := tempWriter.(*os.File); ok {
		logWriter = f
	}
}

// Trace function writers a message about the start and exit of a method. Errors are ignored from Write([]byte) method.
func Trace(msg string) func() {
	if logWriter == nil {
		return nil
	}

	value := fmt.Sprintf("[%s] starting: %s\n", time.Now().Format("15:04:56"), msg)
	logWriter.Write([]byte(value))

	return func() {
		if logWriter == nil {
			return
		}

		value := fmt.Sprintf("[%s] exiting: %s\n", time.Now().Format("15:04:56"), msg)

		logWriter.Write([]byte(value))
	}
}
