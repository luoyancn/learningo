package conf

import (
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once
var overwrite sync.Once

func init() {
	once.Do(func() {
		set_default_section()
		set_kubernetes_section()
		set_etcd_section()
		set_calico_section()
		set_docker_section()
		set_ca_section()
	})
}

func OverWriteConf() {
	overwrite.Do(func() {
		over_write_default_section()
		over_write_kubernetes_section()
		over_write_etcd_section()
		over_write_calico_section()
		over_write_docker_section()
		over_write_ca_section()
	})
}

func set_default_section() {
	viper.SetDefault("default.debug", false)
	viper.SetDefault("default.log_file", "k8sdeploy.log")
}

func over_write_default_section() {
	DEBUG = viper.GetBool("default.debug")
	LOGFILE = viper.GetString("default.log_file")
}

func set_kubernetes_section() {
	viper.SetDefault("kubernetes.k8s_bin_path", "/usr/local/bin")
	viper.SetDefault("kubernetes.k8s_binary", "kubernetes")
	viper.SetDefault("kubernetes.k8s_nodes", map[string]string{})
	viper.SetDefault("kubernetes.overwrite_file", true)
	viper.SetDefault("kubernetes.overwrite_binary", false)
}

func over_write_kubernetes_section() {
	KUBERNETES_K8S_BIN_PATH = viper.GetString("kubernetes.k8s_bin_path")
	KUBERNETES_K8S_BINARY = viper.GetString("kubernetes.k8s_binary")
	KUBERNETES_K8S_NODES = viper.GetStringMapString("kubernetes.k8s_nodes")
	KUBERNETES_OVERWRITE_FILE = viper.GetBool("kubernetes.overwrite_file")
	KUBERNETES_OVERWRITE_BINARY = viper.GetBool("kubernetes.overwrite_binary")
}

func set_etcd_section() {
	viper.SetDefault("etcd.template", "template")
	viper.SetDefault("etcd.nodes", map[string]string{})
	viper.SetDefault("etcd.overwrite", true)
}

func over_write_etcd_section() {
	ETCD_TEMPLATE = viper.GetString("etcd.template")
	ETCD_NODES = viper.GetStringMapString("etcd.nodes")
	ETCD_OVERWRITE = viper.GetBool("etcd.overwrite")
}

func set_calico_section() {
	viper.SetDefault("calico.cni_bin_path", "/usr/local/bin")
	viper.SetDefault("calico.cni_conf_path", "/etc/cni/net.d")
	viper.SetDefault("calico.cni_binary", "")
	viper.SetDefault("calico.overwrite_file", true)
	viper.SetDefault("calico.overwrite_binary", false)
}

func over_write_calico_section() {
	CALICO_CNI_BIN_PATH = viper.GetString("calico.cni_bin_path")
	CALICO_CNI_CONF_PATH = viper.GetString("calico.cni_conf_path")
	CALICO_CNI_BINARY = viper.GetString("calico.cni_binary")
	CALICO_OVERWRITE_FILE = viper.GetBool("calico.overwrite_file")
	CALICO_OVERWRITE_BINARY = viper.GetBool("calico.overwrite_binary")
}

func set_docker_section() {
	viper.SetDefault("docker.template", "docker.service.template")
	viper.SetDefault("docker.docker_hub_mirror",
		"https://docker.mirrors.ustc.edu.cn")
	viper.SetDefault("docker.insecure_registrys", []string{})
	viper.SetDefault("docker.overwrite", true)
}

func over_write_docker_section() {
	DOCKER_TEMPLATE = viper.GetString("docker.template")
	DOCKER_HUB_MIRROR = viper.GetString("docker.docker_hub_mirror")
	DOCKER_INSECURE_REGISTRYS = viper.GetStringSlice(
		"docker.insecure_registrys")
	DOCKER_OVERWRITE = viper.GetBool("docker.overwrite")
}

func set_ca_section() {
	viper.SetDefault("ca.template_path", "initca")
	viper.SetDefault("ca.output", "ca")
	viper.SetDefault("ca.overwrite", true)
}

func over_write_ca_section() {
	CA_TEMPLATE_PATH = viper.GetString("ca.template_path")
	CA_OUTPUT = viper.GetString("ca.output")
	CA_OVERWRITE = viper.GetBool("ca.overwrite")
}
