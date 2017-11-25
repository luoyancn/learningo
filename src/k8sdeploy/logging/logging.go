package logging

import (
	"log"
	"os"

	gologging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

var LOG = gologging.MustGetLogger("k8sdeploy")

var format_std = gologging.MustStringFormatter(
	"%{color}%{time:2006-01-02 15:04:05.999999}" +
		" [%{level:.8s}] %{shortfile} %{shortfunc}" +
		" %{color:reset} %{message}",
)

var format_file = gologging.MustStringFormatter(
	"%{time:2006-01-02 15:04:05.999999} [%{level:.8s}] " +
		" %{shortfile} %{shortfunc} %{message}",
)

func GetLogger() {
	logfile, err := os.OpenFile(
		viper.GetString("default.log_file"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	file_backend := gologging.NewLogBackend(logfile, "", 0)
	std_backend := gologging.NewLogBackend(os.Stdout, "", 0)

	file_back_formater := gologging.NewBackendFormatter(
		file_backend, format_file)
	std_back_formater := gologging.NewBackendFormatter(
		std_backend, format_std)

	gologging.SetBackend(file_back_formater, std_back_formater)
	if !viper.GetBool("default.debug") {
		gologging.SetLevel(gologging.DEBUG, "")
	} else {
		gologging.SetLevel(gologging.INFO, "")
	}
}
