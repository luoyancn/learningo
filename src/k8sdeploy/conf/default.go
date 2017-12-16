package conf

import (
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once

var (

	// Global [defalut] configurations
	DEBUG = false

	// Calico configurations in [calico] section
	CNI_CONF_PATH = "/etc/cni/net.d"

	// Kubelet configurations in [k8s] section
	K8S_NODES               = map[string]string{}
	KUBE_CADVISOR_PORT      = 4194
	KUBE_CLUSTER_DNS_SVC_IP = "10.20.0.2"
	KUBE_DNS_DOMAIN         = "k8s.zhangjl.me"
	KUBELET_HEALTHZ_PORT    = 10248
	KUBE_POD_INFRA_IMAGE    = "gcr.io/google_containers/pause:3.0"
	KUBELET_PORT            = 10250
	KUBELET_READONLY_PORT   = 10255
	SCHEDULER_TEMPLATE      = ""
	KUBELET_TEMPLATE        = ""

	// Ectd configurations in [etcd] section
	ETCD_NODES = map[string]string{}
)

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
	viper.SetDefault("etcd.template", "etcd.conf.template")
	viper.SetDefault("etcd.protocal", "https")
	viper.SetDefault("etcd.ssl", "/etc/etcd/ssl")
	viper.SetDefault("etcd.debug", true)
	viper.SetDefault("etcd.client_cert_auth", true)
	viper.SetDefault("etcd.peer_cert_auth", true)
	viper.SetDefault("etcd.cluster_token", "k8s-etcd-cluster")
	viper.SetDefault("etcd.cluster_name", "k8s-etcd-cluster")
}

func OverWriteConf() {
	once.Do(func() {
		CNI_CONF_PATH = viper.GetString("calico.cni_conf_path")
		K8S_NODES = viper.GetStringMapString("k8s.nodes")
		KUBE_CADVISOR_PORT = viper.GetInt("k8s.kube_cadvisor_port")
		KUBE_CLUSTER_DNS_SVC_IP = viper.GetString(
			"k8s.kube_cluster_dns_svc_ip")
		KUBE_DNS_DOMAIN = viper.GetString("k8s.kube_dns_domain")
		KUBELET_HEALTHZ_PORT = viper.GetInt("k8s.kubelet_healthz_port")
		KUBE_POD_INFRA_IMAGE = viper.GetString("k8s.kube_pod_infra_image")
		KUBELET_PORT = viper.GetInt("k8s.kubelet_port")
		KUBELET_READONLY_PORT = viper.GetInt("k8s.kubelet_readonly_port")
		SCHEDULER_TEMPLATE = viper.GetString("k8s.scheduler_template")
		KUBELET_TEMPLATE = viper.GetString("k8s.kubelet_template")
		ETCD_NODES = viper.GetStringMapString("etcd.nodes")
	})
}
