package logger

import (
	"io"
	"log"
	"os"

	"github.com/influxdata/wlog"
)

// newTelegrafWriter returns a logging-wrapped writer.
func newTelegrafWriter(w io.Writer) io.Writer {
	return &telegrafLog{
		writer: wlog.NewWriter(w),
	}
}

type telegrafLog struct {
	writer io.Writer
}

func (t *telegrafLog) Write(p []byte) (n int, err error) {
	return t.writer.Write(p)
}

// SetupLogging configures the logging output.
//   debug   will set the log level to DEBUG
//   quiet   will set the log level to ERROR
//   logfile will direct the logging output to a file. Empty string is
//           interpreted as stdout. If there is an error opening the file the
//           logger will fallback to stdout.
func SetupLogging(debug, quiet bool, logfile string) {
	if debug {
		wlog.SetLevel(wlog.DEBUG)
	}
	if quiet {
		wlog.SetLevel(wlog.ERROR)
	}

	var oFile *os.File
	if logfile != "" {
		if _, err := os.Stat(logfile); os.IsNotExist(err) {
			if oFile, err = os.Create(logfile); err != nil {
				log.Printf("E! Unable to create %s (%s), using stdout", logfile, err)
				oFile = os.Stdout
			}
		} else {
			if oFile, err = os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY, os.ModeAppend); err != nil {
				log.Printf("E! Unable to append to %s (%s), using stdout", logfile, err)
				oFile = os.Stdout
			}
		}
	} else {
		oFile = os.Stdout
	}

	log.SetOutput(newTelegrafWriter(oFile))
}
