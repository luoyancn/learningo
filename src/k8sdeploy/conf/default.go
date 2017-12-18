package conf

var (

	// Global [defalut] configurations
	DEBUG   = false
	LOGFILE = "k8sdeploy.log"

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
	KUBERNETES_K8S_APISERVER_SECURE_PORT    = 6443
	KUBERNETES_K8S_API_SERVER_TEMPLATE      = ""
	KUBERNETES_K8S_CONTROLLER_TEMPLATE      = ""
	KUBERNETES_K8S_SCHEDULER_TEMPLATE       = ""
	KUBERNETES_K8S_CLUSTER_SERVICE_IP_CIDR  = "10.20.0.0/16"
	KUBERNETES_K8S_APISERVER_INSECURE_PORT  = 8080
	KUBERNETES_K8S_APISERVER_RUNTIME_CONFIG = "rbac.authorization.k8s.io/v1beta1"
	KUBERNETES_K8S_SERVICE_NODE_PORT_RANGE  = "30000-32767"
	KUBERNETES_K8S_CLUSTER_POD_IP_CIDR      = "10.10.0.0/16"
	KUBERNETES_K8S_CONTROLLER_MANAGER_PORT  = 10252
	KUBERNETES_K8S_SCHEDULER_PORT           = 10251

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
