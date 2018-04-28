package common

import (
	"fmt"
	"oceanstack/conf"
	"oceanstack/logging"
	"os"

	"github.com/spf13/viper"
)

func ReadConfig(conf_path string, logger string, logbak int) {
	viper.SetConfigFile(conf_path)
	if err := viper.ReadInConfig(); nil != err {
		fmt.Printf("ERROR:%v\n", err)
		os.Exit(-1)
	}
	conf.OverWriteConf()
	logging.GetLogger(logger, logbak)

	if conf.VERBOSE {
		for key, value := range viper.AllSettings() {
			settings := value.(map[string]interface{})
			for setting_key, setting_value := range settings {
				logging.LOG.Noticef(
					"%s.%s\t%v\n", key, setting_key, setting_value)
			}
		}
	}
}
