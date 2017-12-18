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
	viper.SetDefault("kubernetes.k8s_binary", "")
	viper.SetDefault("kubernetes.k8s_nodes", map[string]string{})
	viper.SetDefault("kubernetes.k8s_overwrite_file", true)
	viper.SetDefault("kubernetes.k8s_overwrite_binary", false)
	viper.SetDefault("kubernetes.k8s_ssl_config_path", "/etc/kubernetes/ssl")
	viper.SetDefault("kubernetes.k8s_config_path", "/etc/kubernetes")
	viper.SetDefault("kubernetes.k8s_cluster_name", "kubernetes")
	viper.SetDefault("kubernetes.k8s_api_server", "")
	viper.SetDefault("kubernetes.k8s_apiserver_secure_port", 6443)
	viper.SetDefault("kubernetes.k8s_api_server_template", "")
	viper.SetDefault("kubernetes.k8s_controller_template", "")
	viper.SetDefault("kubernetes.k8s_scheduler_template", "")
	viper.SetDefault("kubernetes.k8s_cluster_service_ip_cidr", "10.20.0.0/16")
	viper.SetDefault("kubernetes.k8s_apiserver_insecure_port", 8080)
	viper.SetDefault("kubernetes.k8s_apiserver_runtime_config",
		"rbac.authorization.k8s.io/v1beta1")
	viper.SetDefault("kubernetes.k8s_cluster_service_ip_cidr", "10.20.0.0/16")
	viper.SetDefault("kubernetes.k8s_service_node_port_range", "30000-32767")
	viper.SetDefault("kubernetes.k8s_cluster_pod_ip_cidr", "10.10.0.0/16")
	viper.SetDefault("kubernetes.k8s_controller_manager_port", 10252)
	viper.SetDefault("kubernetes.k8s_scheduler_port", 10251)
}

func over_write_kubernetes_section() {
	KUBERNETES_K8S_BIN_PATH = viper.GetString("kubernetes.k8s_bin_path")
	KUBERNETES_K8S_BINARY = viper.GetString("kubernetes.k8s_binary")
	KUBERNETES_K8S_NODES = viper.GetStringMapString("kubernetes.k8s_nodes")
	KUBERNETES_K8S_OVERWRITE_FILE = viper.GetBool(
		"kubernetes.k8s_overwrite_file")
	KUBERNETES_K8S_OVERWRITE_BINARY = viper.GetBool(
		"kubernetes.k8s_overwrite_binary")
	KUBERNETES_K8S_SSL_CONFIG_PATH = viper.GetString(
		"kubernetes.k8s_ssl_config_path")
	KUBERNETES_K8S_CONFIG_PATH = viper.GetString("kubernetes.k8s_config_path")
	KUBERNETES_K8S_CLUSTER_NAME = viper.GetString("kubernetes.k8s_cluster_name")
	KUBERNETES_K8S_API_SERVER = viper.GetString("kubernetes.k8s_api_server")
	KUBERNETES_K8S_APISERVER_SECURE_PORT = viper.GetInt(
		"kubernetes.k8s_apiserver_secure_port")
	KUBERNETES_K8S_API_SERVER_TEMPLATE = viper.GetString(
		"kubernetes.k8s_api_server_template")
	KUBERNETES_K8S_CONTROLLER_TEMPLATE = viper.GetString(
		"kubernetes.k8s_controller_template")
	KUBERNETES_K8S_SCHEDULER_TEMPLATE = viper.GetString(
		"kubernetes.k8s_scheduler_template")
	KUBERNETES_K8S_CLUSTER_SERVICE_IP_CIDR = viper.GetString(
		"kubernetes.k8s_cluster_service_ip_cidr")
	KUBERNETES_K8S_APISERVER_INSECURE_PORT = viper.GetInt(
		"kubernetes.k8s_apiserver_insecure_port")
	KUBERNETES_K8S_APISERVER_RUNTIME_CONFIG = viper.GetString(
		"kubernetes.k8s_apiserver_runtime_config")
	KUBERNETES_K8S_CLUSTER_SERVICE_IP_CIDR = viper.GetString(
		"kubernetes.k8s_cluster_service_ip_cidr")
	KUBERNETES_K8S_SERVICE_NODE_PORT_RANGE = viper.GetString(
		"kubernetes.k8s_service_node_port_range")
	KUBERNETES_K8S_CLUSTER_POD_IP_CIDR = viper.GetString(
		"kubernetes.k8s_cluster_pod_ip_cidr")
	KUBERNETES_K8S_CONTROLLER_MANAGER_PORT = viper.GetInt(
		"kubernetes.k8s_controller_manager_port")
	KUBERNETES_K8S_SCHEDULER_PORT = viper.GetInt(
		"kubernetes.k8s_scheduler_port")
}

func set_etcd_section() {
	viper.SetDefault("etcd.template", "")
	viper.SetDefault("etcd.nodes", map[string]string{})
	viper.SetDefault("etcd.overwrite", true)
	viper.SetDefault("etcd.protocal", "https")
	viper.SetDefault("etcd.token", "k8s-cluster-token")
	viper.SetDefault("etcd.client_cert_auth", true)
	viper.SetDefault("etcd.peer_cert_auth", true)
	viper.SetDefault("etcd.etcd_ssl", "/etc/etcd/ssl")
	viper.SetDefault("etcd.etcd_debug", true)
}

func over_write_etcd_section() {
	ETCD_TEMPLATE = viper.GetString("etcd.template")
	ETCD_NODES = viper.GetStringMapString("etcd.nodes")
	ETCD_OVERWRITE = viper.GetBool("etcd.overwrite")
	ETCD_PROTOCAL = viper.GetString("etcd.protocal")
	ETCD_TOKEN = viper.GetString("etcd.token")
	ETCD_CLIENT_CERT_AUTH = viper.GetBool("etcd.client_cert_auth")
	ETCD_PEER_CERT_AUTH = viper.GetBool("etcd.peer_cert_auth")
	ETCD_SSL = viper.GetString("etcd.etcd_ssl")
	ETCD_DEBUG = viper.GetBool("etcd.etcd_debug")
}

func set_calico_section() {
	viper.SetDefault("calico.cni_bin_path", "/usr/local/bin")
	viper.SetDefault("calico.cni_conf_path", "/etc/cni/net.d")
	viper.SetDefault("calico.cni_binary", "")
	viper.SetDefault("calico.cni_overwrite_file", true)
	viper.SetDefault("calico.cni_overwrite_binary", false)
}

func over_write_calico_section() {
	CALICO_CNI_BIN_PATH = viper.GetString("calico.cni_bin_path")
	CALICO_CNI_CONF_PATH = viper.GetString("calico.cni_conf_path")
	CALICO_CNI_BINARY = viper.GetString("calico.cni_binary")
	CALICO_CNI_OVERWRITE_FILE = viper.GetBool("calico.cni_overwrite_file")
	CALICO_CNI_OVERWRITE_BINARY = viper.GetBool("calico.cni_overwrite_binary")
}

func set_docker_section() {
	viper.SetDefault("docker.template", "")
	viper.SetDefault("docker.docker_hub_mirror", "")
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
	viper.SetDefault("ca.template_path", "")
	viper.SetDefault("ca.output", "")
	viper.SetDefault("ca.overwrite", true)
}

func over_write_ca_section() {
	CA_TEMPLATE_PATH = viper.GetString("ca.template_path")
	CA_OUTPUT = viper.GetString("ca.output")
	CA_OVERWRITE = viper.GetBool("ca.overwrite")
}
