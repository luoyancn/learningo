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
	TRACE = log.New(traceHandle, "TRACE: ",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	INFO = log.New(infoHandle, "INFO: ",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	WARNING = log.New(warningHandle, "WARNING: ",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	ERROR = log.New(errorHandle, "ERROR: ",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	DEBUG = log.New(debugHandler, "DEBUG: ",
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
}

func GetLogger() {
	logfile, err := os.OpenFile(
		viper.GetString("default.log_file"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	multi_writer := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(multi_writer)
	if viper.GetBool("default.debug") {
		initLog(multi_writer, multi_writer, multi_writer,
			multi_writer, multi_writer)
	} else {
		initLog(multi_writer, multi_writer, multi_writer,
			multi_writer, ioutil.Discard)
	}
}
