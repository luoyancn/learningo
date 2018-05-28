package conf

import (
	"runtime"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var once sync.Once
var overwrite sync.Once

func init() {
	once.Do(func() {
		set_default_section()
		set_database_section()
		set_redis_section()
		set_kubernetes_section()
		set_etcd_section()
		set_calico_section()
		set_docker_section()
		set_ca_section()
		set_grpc_section()
		set_api_section()
	})
}

func OverWriteConf() {
	overwrite.Do(func() {
		over_write_default_section()
		over_write_database_section()
		over_write_redis_section()
		over_write_kubernetes_section()
		over_write_etcd_section()
		over_write_calico_section()
		over_write_docker_section()
		over_write_ca_section()
		over_write_grpc_section()
		over_write_api_section()
	})
}

func set_default_section() {
	viper.SetDefault("default.debug", false)
	viper.SetDefault("default.verbose", false)
	viper.SetDefault("default.log_path", "/var/log/oceanstack")
	viper.SetDefault("default.admin_token", "ADMIN_TOKEN")
	viper.SetDefault("default.ssh_port", 22)
	viper.SetDefault("default.ssh_timeout", 5)
}

func over_write_default_section() {
	DEBUG = viper.GetBool("default.debug")
	VERBOSE = viper.GetBool("default.verbose")
	LOGPATH = viper.GetString("default.log_path")
	ADMIN_TOKEN = viper.GetString("default.admin_token")
	SSH_PORT = viper.GetInt("default.ssh_port")
	SSH_TIMEOUT = viper.GetDuration("default.ssh_timeout") * time.Second
}

func set_database_section() {
	viper.SetDefault("database.connection",
		"golang:golang@tcp(127.0.0.1:3306)/golang?parseTime=true&loc=Local")
	viper.SetDefault("database.max_time_min", 30)
	viper.SetDefault("database.max_idle", 30)
	viper.SetDefault("database.max_open", 30)
	viper.SetDefault("database.debug_mode", false)
}

func over_write_database_section() {
	DATABASE_CONNECTION = viper.GetString("database.connection")
	DATABASE_MAX_TIME_MIN = viper.GetDuration(
		"database.max_time_min") * time.Minute
	DATABASE_MAX_IDLE = viper.GetInt("database.max_idle")
	DATABASE_MAX_OPEN = viper.GetInt("database.max_open")
	DATABASE_DEBUG_MODE = viper.GetBool("database.debug_mode")
}

func set_redis_section() {
	viper.SetDefault("redis.connection", "127.0.0.1:6379")
	viper.SetDefault("redis.max_idle", 30)
	viper.SetDefault("redis.max_active", 30)
	viper.SetDefault("redis.max_conn_lifetime", 10)
	viper.SetDefault("redis.idle_timeout", 10)
	viper.SetDefault("redis.database", 0)
	viper.SetDefault("redis.expire", 30*60)
}

func over_write_redis_section() {
	REDIS_CONNECTION = viper.GetString("redis.connection")
	REDIS_MAX_IDLE = viper.GetInt("redis.max_idle")
	REDIS_MAX_ACTIVE = viper.GetInt("redis.max_active")
	REDIS_MAX_CONN_LIFETIME = viper.GetDuration(
		"redis.max_conn_lifetime") * time.Minute
	REDIS_IDLE_TIMEOUT = viper.GetDuration(
		"redis.idle_timeout") * time.Minute
	REDIS_DATABASE = viper.GetInt("redis.database")
	REDIS_EXPIRE = viper.GetInt("redis.expire") * 60
}

func set_grpc_section() {
	viper.SetDefault("grpc.server", "127.0.0.1")
	viper.SetDefault("grpc.port", 36000)
	viper.SetDefault("grpc.workers", runtime.NumCPU()-1)
	viper.SetDefault("grpc.timeout", 5)
	viper.SetDefault("grpc.pool_size", 10)
	viper.SetDefault("grpc.concurrency", 64)
	viper.SetDefault("grpc.server_conn_limits", 1024)
	viper.SetDefault("grpc.server_req_max_frequency", 1024)
	viper.SetDefault("grpc.server_req_burst_frequency", 10)
	viper.SetDefault("grpc.req_msg_size", 1)
	viper.SetDefault("grpc.enable_lb", false)
	viper.SetDefault("grpc.lb_listen", "127.0.0.1:9000")
	viper.SetDefault("grpc.etcd_endpoints", []string{"http//:127.0.0.1:2379"})
	viper.SetDefault("grpc.service_name", "oceanstack")
	viper.SetDefault("grpc.etcd_ca", "/etc/etcd/ssl/ca.pem")
	viper.SetDefault("grpc.etcd_cert", "/etc/etcd/ssl/cert.pem")
	viper.SetDefault("grpc.etcd_key", "/etc/etcd/ssl/key.pem")
	viper.SetDefault("grpc.etcd_timeout", 30)
}

func over_write_grpc_section() {
	GRPC_SERVER = viper.GetString("grpc.server")
	GRPC_PORT = viper.GetInt("grpc.port")
	GRPC_WORKERS = viper.GetInt("grpc.workers")
	GRPC_TIMEOUT = viper.GetDuration("grpc.timeout") * time.Second
	GRPC_POOL_SIZE = viper.GetInt("grpc.pool_size")
	GRPC_CONCURRENCY = viper.GetInt("grpc.concurrency")
	GRPC_SERVER_CONN_LIMITS = viper.GetInt("grpc.server_conn_limits")
	GRPC_SERVER_REQ_MAX_FREQUENCY = viper.GetFloat64(
		"grpc.server_req_max_frequency")
	GRPC_SERVER_REQ_BURST_FREQUENCY = viper.GetInt(
		"grpc.server_req_burst_frequency")
	GRPC_REQ_MSG_SIZE = viper.GetInt("grpc.req_msg_size") * 1024 * 1024
	GRPC_ENABLE_LB = viper.GetBool("grpc.enable_lb")
	GRPC_LB_LISTEN = viper.GetString("grpc.lb_listen")
	GRPC_ETCD_ENDPOINTS = viper.GetStringSlice("grpc.etcd_endpoints")
	GRPC_ETCD_SERVICE_NAME = viper.GetString("grpc.service_name")
	GRPC_ETCD_CA = viper.GetString("grpc.etcd_ca")
	GRPC_ETCD_CERT = viper.GetString("grpc.etcd_cert")
	GRPC_ETCD_KEY = viper.GetString("grpc.etcd_key")
	GRPC_ETCD_TIMEOUT = viper.GetDuration("grpc.etcd_timeout") * time.Second
}

func set_api_section() {
	viper.SetDefault("api.listen", "127.0.0.1:8888")
	viper.SetDefault("api.workers", runtime.NumCPU()-1)
}

func over_write_api_section() {
	API_LISTEN = viper.GetString("api.listen")
	API_WORKERS = viper.GetInt("api.workers")
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
	viper.SetDefault("kubernetes.k8s_apiserver_secure_port", 5443)
	viper.SetDefault("kubernetes.k8s_api_server_template", "")
	viper.SetDefault("kubernetes.k8s_controller_template", "")
	viper.SetDefault("kubernetes.k8s_scheduler_template", "")
	viper.SetDefault("kubernetes.k8s_apiserver_insecure_port", 7070)
	viper.SetDefault("kubernetes.k8s_apiserver_runtime_config",
		"rbac.authorization.k8s.io/v1beta1")
	viper.SetDefault("kubernetes.k8s_cluster_service_ip_cidr", "10.20.0.0/16")
	viper.SetDefault("kubernetes.k8s_service_node_port_range", "30000-32767")
	viper.SetDefault("kubernetes.k8s_cluster_pod_ip_cidr", "10.10.0.0/16")
	viper.SetDefault("kubernetes.k8s_controller_manager_port", 10252)
	viper.SetDefault("kubernetes.k8s_scheduler_port", 10251)
	viper.SetDefault("kubernetes.kubelete_template", "")
	viper.SetDefault("kubernetes.kube_cadvisor_port", 4194)
	viper.SetDefault("kubernetes.kube_cluster_dns_svc_ip", "10.20.0.2")
	viper.SetDefault("kubernetes.kube_dns_domain", "k8s.zhangjl.me")
	viper.SetDefault("kubernetes.kubelet_healthz_port", 10248)
	viper.SetDefault("kubernetes.kube_pod_infra_image",
		"gcr.io/google_containers/pause:3.0")
	viper.SetDefault("kubernetes.kubelet_port", 10250)
	viper.SetDefault("kubernetes.kubelet_readonly_port", 10255)
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
	KUBERNETES_KUBELET_TEMPLATE = viper.GetString(
		"kubernetes.kubelet_template")
	KUBERNETES_KUBE_CADVISOR_PORT = viper.GetInt(
		"kubernetes.kube_cadvisor_port")
	KUBERNETES_KUBE_CLUSTER_DNS_SVC_IP = viper.GetString(
		"kubernetes.kube_cluster_dns_svc_ip")
	KUBERNETES_KUBE_DNS_DOMAIN = viper.GetString("kubernetes.kube_dns_domain")
	KUBERNETES_KUBELET_HEALTHZ_PORT = viper.GetInt(
		"kubernetes.kubelet_healthz_port")
	KUBERNETES_KUBE_POD_INFRA_IMAGE = viper.GetString(
		"kubernetes.kube_pod_infra_image")
	KUBERNETES_KUBELET_PORT = viper.GetInt("kubernetes.kubelet_port")
	KUBERNETES_KUBELET_READONLY_PORT = viper.GetInt(
		"kubernetes.kubelet_readonly_port")
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
