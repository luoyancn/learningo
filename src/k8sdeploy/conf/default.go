package conf

var (

	// Global [defalut] configurations
	DEBUG   = false
	LOGFILE = "k8sdeploy.log"

	// Ca configurations in [ca] section
	CA_TEMPLATE_PATH = "initca"
	CA_OUTPUT        = "ca"
	CA_OVERWRITE     = true

	// Calico configurations in [calico] section
	CALICO_CNI_BIN_PATH     = "/usr/local/bin"
	CALICO_CNI_CONF_PATH    = "/etc/cni/net.d"
	CALICO_CNI_BINARY       = "/home/zhangjl/kubernetes/calico"
	CALICO_OVERWRITE_FILE   = true
	CALICO_OVERWRITE_BINARY = false

	// kubernetes configurations in [kubernetes] section
	KUBERNETES_K8S_BIN_PATH     = "/usr/local/bin"
	KUBERNETES_K8S_BINARY       = "kubernetes"
	KUBERNETES_K8S_NODES        = map[string]string{}
	KUBERNETES_OVERWRITE_FILE   = true
	KUBERNETES_OVERWRITE_BINARY = false

	// Ectd configurations in [etcd] section
	// And, etcd server always means kubernetes nodes, but etcd use private
	// ip address
	ETCD_TEMPLATE  = ""
	ETCD_NODES     = map[string]string{}
	ETCD_OVERWRITE = true

	// Docker configurations in [docker] section
	DOCKER_TEMPLATE           = "docker.service.template"
	DOCKER_HUB_MIRROR         = "https://docker.mirrors.ustc.edu.cn"
	DOCKER_INSECURE_REGISTRYS = []string{}
	DOCKER_OVERWRITE          = true
)
