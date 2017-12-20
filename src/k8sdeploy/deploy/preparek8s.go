package deploy

import (
	"k8sdeploy/conf"
	"k8sdeploy/utils"
)

func PrepareK8SBinary() bool {
	nodes := []string{}
	for _, node := range conf.KUBERNETES_K8S_NODES {
		nodes = append(nodes, node)
	}
	binary_files_path := conf.KUBERNETES_K8S_BINARY
	dest_binary_path := conf.KUBERNETES_K8S_BIN_PATH
	overwrite_k8s_binary := conf.KUBERNETES_K8S_OVERWRITE_BINARY
	return utils.SCPFiles([]string{binary_files_path},
		dest_binary_path, "binary", overwrite_k8s_binary, nodes...)
}

func PrepareCAKey() bool {
	nodes := []string{}
	for _, node := range conf.KUBERNETES_K8S_NODES {
		nodes = append(nodes, node)
	}
	source_json_path := conf.CA_TEMPLATE_PATH
	source_ca_path := conf.CA_OUTPUT
	dest_ca_path := conf.KUBERNETES_K8S_SSL_CONFIG_PATH
	overwrite_ca := conf.CA_OVERWRITE
	return utils.SCPFiles([]string{source_json_path, source_ca_path},
		dest_ca_path, "", overwrite_ca, nodes...)
}
