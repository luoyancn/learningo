package logging

import (
	"oceanstack/conf"
	"os"
	"path"
	"sync"

	gologging "github.com/op/go-logging"
)

var LOG *gologging.Logger

var format_std = gologging.MustStringFormatter(
	"%{color}%{time:2006-01-02 15:04:05.999}" +
		" [%{level:.8s}] %{shortfile} %{shortfunc}" +
		" %{color:reset} %{message}",
)

var format_file = gologging.MustStringFormatter(
	"%{time:2006-01-02 15:04:05.999} [%{level:.8s}] " +
		" [%{shortfile}] [%{shortfunc}] %{message}",
)

var once sync.Once

const STD_ENABLED = 1
const FILE_ENABLED = 2

func GetLogger(logger string, logback int) {
	once.Do(func() {
		LOG = gologging.MustGetLogger(logger)
		logfile, err := os.OpenFile(
			path.Join(conf.LOGPATH, logger),
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			LOG.Panicf("Cannot create the log file :%v\n", err)
		}
		file_backend := gologging.NewLogBackend(logfile, "", 0)
		file_back_formater := gologging.NewBackendFormatter(
			file_backend, format_file)
		std_backend := gologging.NewLogBackend(os.Stdout, "", 0)
		std_back_formater := gologging.NewBackendFormatter(
			std_backend, format_std)
		switch logback {
		case 2:
			gologging.SetBackend(file_back_formater)
			break
		case 3:
			gologging.SetBackend(file_back_formater, std_back_formater)
			break
		default:
			gologging.SetBackend(std_back_formater)
			break
		}

		if conf.DEBUG {
			gologging.SetLevel(gologging.DEBUG, "")
		} else {
			gologging.SetLevel(gologging.INFO, "")
		}
	})
}
