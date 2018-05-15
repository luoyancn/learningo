package conf

import (
	"time"
)

var (

	// Global [defalut] configurations
	DEBUG       = false
	VERBOSE     = false
	LOGPATH     = "/var/log/oceanstack"
	ADMIN_TOKEN = "ADMIN_TOKEN"
	SSH_PORT    = 22
	SSH_TIMEOUT = 5 * time.Second
	LISTEN      = "127.0.0.1:8888"

	// Global [database] configurations
	DATABASE_CONNECTION   = "golang:golang@tcp(127.0.0.1:3306)/golang?parseTime=true&loc=Local"
	DATABASE_MAX_TIME_MIN = 30 * time.Minute
	DATABASE_MAX_IDLE     = 30
	DATABASE_MAX_OPEN     = 30
	DATABASE_DEBUG_MODE   = false

	// Global [redis] configurations
	REDIS_CONNECTION        = "127.0.0.1:6379"
	REDIS_MAX_IDLE          = 30
	REDIS_MAX_ACTIVE        = 30
	REDIS_MAX_CONN_LIFETIME = 10 * time.Minute
	REDIS_IDLE_TIMEOUT      = 10 * time.Minute
	REDIS_DATABASE          = 0
	REDIS_EXPIRE            = 30 * 60

	// Ca configurations in [ca] section
	CA_TEMPLATE_PATH = ""
	CA_OUTPUT        = ""
	CA_OVERWRITE     = true

	// Calico configurations in [calico] section
	CALICO_CNI_BIN_PATH         = "/usr/local/bin"
	CALICO_CNI_CONF_PATH        = "/etc/cni/net.d"
	CALICO_CNI_BINARY           = ""
	CALICO_CNI_OVERWRITE_FILE   = true
	CALICO_CNI_OVERWRITE_BINARY = false

	// kubernetes configurations in [kubernetes] section
	KUBERNETES_K8S_BIN_PATH                 = "/usr/local/bin"
	KUBERNETES_K8S_BINARY                   = ""
	KUBERNETES_K8S_NODES                    = map[string]string{}
	KUBERNETES_K8S_OVERWRITE_FILE           = true
	KUBERNETES_K8S_OVERWRITE_BINARY         = false
	KUBERNETES_K8S_SSL_CONFIG_PATH          = "/etc/kubernetes/ssl"
	KUBERNETES_K8S_CONFIG_PATH              = "/etc/kubernetes"
	KUBERNETES_K8S_CLUSTER_NAME             = "kubernetes"
	KUBERNETES_K8S_API_SERVER               = ""
	KUBERNETES_K8S_APISERVER_SECURE_PORT    = 5443
	KUBERNETES_K8S_API_SERVER_TEMPLATE      = ""
	KUBERNETES_K8S_CONTROLLER_TEMPLATE      = ""
	KUBERNETES_K8S_SCHEDULER_TEMPLATE       = ""
	KUBERNETES_K8S_CLUSTER_SERVICE_IP_CIDR  = "10.20.0.0/16"
	KUBERNETES_K8S_APISERVER_INSECURE_PORT  = 7070
	KUBERNETES_K8S_APISERVER_RUNTIME_CONFIG = "rbac.authorization.k8s.io/v1beta1"
	KUBERNETES_K8S_SERVICE_NODE_PORT_RANGE  = "30000-32767"
	KUBERNETES_K8S_CLUSTER_POD_IP_CIDR      = "10.10.0.0/16"
	KUBERNETES_K8S_CONTROLLER_MANAGER_PORT  = 10252
	KUBERNETES_K8S_SCHEDULER_PORT           = 10251
	KUBERNETES_KUBELET_TEMPLATE             = ""
	KUBERNETES_KUBE_CADVISOR_PORT           = 4194
	KUBERNETES_KUBE_CLUSTER_DNS_SVC_IP      = "10.20.0.2"
	KUBERNETES_KUBE_DNS_DOMAIN              = "k8s.zhangjl.me"
	KUBERNETES_KUBELET_HEALTHZ_PORT         = 10248
	KUBERNETES_KUBE_POD_INFRA_IMAGE         = "gcr.io/google_containers/pause:3.0"
	KUBERNETES_KUBELET_PORT                 = 10250
	KUBERNETES_KUBELET_READONLY_PORT        = 10255
	// Ectd configurations in [etcd] section
	// And, etcd server always means kubernetes nodes, but etcd use private
	// ip address
	ETCD_TEMPLATE         = ""
	ETCD_NODES            = map[string]string{}
	ETCD_OVERWRITE        = true
	ETCD_PROTOCAL         = "https"
	ETCD_TOKEN            = "k8s-cluster-token"
	ETCD_CLIENT_CERT_AUTH = true
	ETCD_PEER_CERT_AUTH   = true
	ETCD_SSL              = "/etc/etcd/ssl"
	ETCD_DEBUG            = true

	// Docker configurations in [docker] section
	DOCKER_TEMPLATE           = ""
	DOCKER_HUB_MIRROR         = ""
	DOCKER_INSECURE_REGISTRYS = []string{}
	DOCKER_OVERWRITE          = true
)
