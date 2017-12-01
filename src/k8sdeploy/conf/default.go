package conf

import (
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once

func init() {
	once.Do(func() {
		set_default_section()
		set_k8s_section()
		set_cfs_section()
	})
}

func set_default_section() {
	viper.SetDefault("default.debug", false)
	viper.SetDefault("default.log_file", "k8sdeploy.log")
}

func set_k8s_section() {
	viper.SetDefault("k8s.target_path", "/usr/bin")
	viper.SetDefault("k8s.config_path", "/etc/kubernetes")
	viper.SetDefault("k8s.ssl_config_path", "/etc/kubernetes/ssl")
}

func set_cfs_section() {
	viper.SetDefault("cfs.output", "ca")
}
