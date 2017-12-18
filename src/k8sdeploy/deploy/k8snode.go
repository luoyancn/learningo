package deploy

import (
	"k8sdeploy/conf"
	"k8sdeploy/logging"
)

func DeployK8sNode() {
	logging.LOG.Noticef("%v\n", conf.CALICO_CNI_BIN_PATH)
	/*
		logging.LOG.Noticef("%v\n", conf.K8S_NODES)
			kubelet_ctx := template.Must(template.ParseFiles(conf.KUBELET_TEMPLATE))
			kubelet_map := map[string]interface{}{
				"kube_cadvisor_port":      conf.KUBE_CADVISOR_PORT,
				"kube_cluster_dns_svc_ip": conf.KUBE_CLUSTER_DNS_SVC_IP,
				"kube_dns_domain":         conf.KUBE_DNS_DOMAIN,
				"kubelet_healthz_port":    conf.KUBELET_HEALTHZ_PORT,
				"kube_pod_infra_image":    conf.KUBE_POD_INFRA_IMAGE,
				"kubelet_port":            conf.KUBELET_PORT,
				"kubelet_readonly_port":   conf.KUBELET_READONLY_PORT}
	*/
}
