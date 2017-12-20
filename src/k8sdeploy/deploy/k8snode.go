package deploy

import (
	"k8sdeploy/conf"
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"text/template"
)

func DeployK8sNode() bool {
	kubelet_ctx := template.Must(
		template.ParseFiles(conf.KUBERNETES_KUBELET_TEMPLATE))
	kubelet_map := map[string]interface{}{
		"kube_cadvisor_port":      conf.KUBERNETES_KUBE_CADVISOR_PORT,
		"kube_cluster_dns_svc_ip": conf.KUBERNETES_KUBE_CLUSTER_DNS_SVC_IP,
		"kube_dns_domain":         conf.KUBERNETES_KUBE_DNS_DOMAIN,
		"kubelet_healthz_port":    conf.KUBERNETES_KUBELET_HEALTHZ_PORT,
		"kube_pod_infra_image":    conf.KUBERNETES_KUBE_POD_INFRA_IMAGE,
		"kubelet_port":            conf.KUBERNETES_KUBELET_PORT,
		"kubelet_readonly_port":   conf.KUBERNETES_KUBELET_READONLY_PORT}

	ips := []string{}
	for host, ip := range conf.ETCD_NODES {
		kubelet_writer, err := os.Create("/tmp/kubelet.service." + host)
		if nil != err {
			logging.LOG.Errorf(
				"Cannot create kubelet service config file:%v\n", err)
			return false
		}
		kubelet_map["node_ip"] = ip
		if err = kubelet_ctx.Execute(
			kubelet_writer, kubelet_map); nil != err {
			logging.LOG.Errorf(
				"Cannot parse kubelet config file:%v\n", err)
			return false
		}
		if !utils.SCPFiles([]string{kubelet_writer.Name()},
			"/usr/lib/systemd/system/kubelet.service",
			"file", true, conf.KUBERNETES_K8S_NODES[host]) {
			return false
		}
		ips = append(ips, conf.KUBERNETES_K8S_NODES[host])
	}
	cmd := "mkdir -p /var/lib/kubelet;systemctl daemon-reload;systemctl enable kubelet;systemctl restart kubelet"
	return utils.RemoteCmd(cmd, ips...)
}
