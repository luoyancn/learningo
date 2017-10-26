package conf

import (
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once

func init() {
	once.Do(func() {
		set_default_section()
		set_database_section()
	})
}

func set_default_section() {
	viper.SetDefault("default.listen", "127.0.0.1:8080")
	viper.SetDefault("default.debug", false)
	viper.SetDefault("default.log_file", "rest.log")
}

func set_database_section() {
	viper.SetDefault("database.debug", false)
	viper.SetDefault("database.connection",
		"golang:golang@tcp(127.0.0.1:3306)/golang?parseTime=true")
	viper.SetDefault("database.max_time", 30)
	viper.SetDefault("database.max_idle", 30)
	viper.SetDefault("database.max_open", 90)
}
