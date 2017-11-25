package main

import (
	"os"

	gologging "github.com/op/go-logging"
)

var log = gologging.MustGetLogger("example")

var format_stdio = gologging.MustStringFormatter(
	"%{color}%{time:2006-01-02 15:04:05.999999} " +
		"%{shortfile} %{shortfunc} [%{level:.4s}] %{color:reset} %{message}",
)
var format_file = gologging.MustStringFormatter(
	`%{time:15:04:05.000} %{shortfunc} [%{level:.4s}] %{message}`,
)

type Password string

func (p Password) Redacted() interface{} {
	return gologging.Redact(string(p))
}

func main() {
	logfile, err := os.OpenFile("log_file",
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(-1)
	}
	backend1 := gologging.NewLogBackend(logfile, "", 0)
	backend2 := gologging.NewLogBackend(os.Stderr, "", 0)

	backend1Formatter := gologging.NewBackendFormatter(backend1, format_file)
	backend2Formatter := gologging.NewBackendFormatter(backend2, format_stdio)

	gologging.SetBackend(backend1Formatter, backend2Formatter)
	gologging.SetLevel(gologging.DEBUG, "")

	log.Debugf("debug %s", Password("secret"))
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("err")
	log.Critical("crit")
}
