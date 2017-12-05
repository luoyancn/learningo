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
		set_etcd_section()
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
	viper.SetDefault("k8s.overwrite_binary", false)
	viper.SetDefault("k8s.overwrite_ssl", true)
	viper.SetDefault("k8s.cluster_name", "kubernetes")
	viper.SetDefault("k8s.tmp", "/tmp")
}

func set_cfs_section() {
	viper.SetDefault("cfs.output", "ca")
}

func set_etcd_section() {
	viper.SetDefault("etcd.workdir_prefix", "/var/lib/etcd/")
	viper.SetDefault("etcd.protocal", "https")
	viper.SetDefault("etcd.cluster_name", "k8s-ectd-cluster")
}
