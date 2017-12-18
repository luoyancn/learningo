package logging

import (
	"k8sdeploy/conf"
	"os"
	"sync"

	gologging "github.com/op/go-logging"
)

var LOG *gologging.Logger

var format_std = gologging.MustStringFormatter(
	"%{color}%{time:2006-01-02 15:04:05.999999999}" +
		" [%{level:.8s}] %{shortfile} %{shortfunc}" +
		" %{color:reset} %{message}",
)

var format_file = gologging.MustStringFormatter(
	"%{time:2006-01-02 15:04:05.999999999} [%{level:.8s}] " +
		" %{shortfile} %{shortfunc} %{message}",
)

var once sync.Once

func GetLogger() {
	once.Do(func() {
		LOG = gologging.MustGetLogger("k8sdeploy")
		logfile, err := os.OpenFile(
			conf.LOGFILE,
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			LOG.Panicf("Cannot create the log file :%v\n", err)
		}
		file_backend := gologging.NewLogBackend(logfile, "", 0)
		std_backend := gologging.NewLogBackend(os.Stdout, "", 0)

		file_back_formater := gologging.NewBackendFormatter(
			file_backend, format_file)
		std_back_formater := gologging.NewBackendFormatter(
			std_backend, format_std)

		gologging.SetBackend(file_back_formater, std_back_formater)
		if conf.DEBUG {
			gologging.SetLevel(gologging.DEBUG, "")
		} else {
			gologging.SetLevel(gologging.INFO, "")
		}
	})
}
