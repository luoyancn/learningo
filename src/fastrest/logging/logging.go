package logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	TRACE   *log.Logger
	INFO    *log.Logger
	WARNING *log.Logger
	ERROR   *log.Logger
	DEBUG   *log.Logger
)

func initLog(traceHandle io.Writer, infoHandle io.Writer,
	warningHandle io.Writer, errorHandle io.Writer, debugHandler io.Writer) {
	TRACE = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	INFO = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WARNING = log.New(warningHandle, "WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DEBUG = log.New(debugHandler, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func GetLogger() {
	logfile, err := os.OpenFile(
		viper.GetString("default.log_file"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("%v\n", err)
	}
	if viper.GetBool("default.debug") {
		initLog(logfile, logfile, logfile, logfile, logfile)
	} else {
		initLog(logfile, logfile, logfile, logfile, ioutil.Discard)
	}
}
