package logging

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	trace   *log.Logger
	info    *log.Logger
	warning *log.Logger
	eRROR   *log.Logger
	debug   *log.Logger
)

func initLog(traceHandle io.Writer, infoHandle io.Writer,
	warningHandle io.Writer, errorHandle io.Writer, debugHandler io.Writer) {
	trace = log.New(
		traceHandle, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	info = log.New(
		infoHandle, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	warning = log.New(
		warningHandle, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	eRROR = log.New(
		errorHandle, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	debug = log.New(
		debugHandler, "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func Trace(format string, msg ...interface{}) {
	trace.Printf("[TRACE] %s", fmt.Sprintf(format, msg...))
}

func Info(format string, msg ...interface{}) {
	info.Printf("[INFO] %s", fmt.Sprintf(format, msg...))
}

func Warning(format string, msg ...interface{}) {
	warning.Printf("[WARNING] %s", fmt.Sprintf(format, msg...))
}

func Error(format string, msg ...interface{}) {
	eRROR.Printf("[ERROR] %s", fmt.Sprintf(format, msg...))
}

func Debug(format string, msg ...interface{}) {
	debug.Printf("[DEBUG] %s", fmt.Sprintf(format, msg...))
}

func GetLogger() {
	logfile, err := os.OpenFile(
		viper.GetString("default.log_file"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	multi_writer := io.MultiWriter(os.Stdout, logfile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	log.SetOutput(multi_writer)
	if viper.GetBool("default.debug") {
		initLog(multi_writer, multi_writer, multi_writer,
			multi_writer, multi_writer)
	} else {
		initLog(multi_writer, multi_writer, multi_writer,
			multi_writer, ioutil.Discard)
	}
}
