package deploy

import (
	"fmt"
	"k8sdeploy/conf"
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"strings"
	"text/template"
)

func DeployEtcd() bool {
	ctx := template.Must(template.ParseFiles(
		conf.ETCD_TEMPLATE))

	map_ctx := map[string]interface{}{
		"etcd_protocal":    conf.ETCD_PROTOCAL,
		"etcd_token":       conf.ETCD_TOKEN,
		"client_cert_auth": conf.ETCD_CLIENT_CERT_AUTH,
		"peer_cert_auth":   conf.ETCD_PEER_CERT_AUTH,
		"etcd_ssl":         conf.ETCD_SSL,
		"etcd_debug":       conf.ETCD_DEBUG}

	etcd_nodes := conf.ETCD_NODES
	nodes := []string{}
	ips := []string{}
	endpoints := []string{}
	for name, ip := range etcd_nodes {
		nodes = append(
			nodes, name+"="+conf.ETCD_PROTOCAL+"://"+ip+":2380")
		endpoints = append(
			endpoints, conf.ETCD_PROTOCAL+"://"+ip+":2379")
	}
	etcd_cluster := strings.Join(nodes, ",")
	etcd_endpoints := strings.Join(endpoints, ",")
	map_ctx["etcd_cluster"] = etcd_cluster
	for name, ip := range etcd_nodes {
		map_ctx["etcd_name"] = name
		map_ctx["etcd_node_ip"] = ip
		writer, err := os.Create("/tmp/" + map_ctx["etcd_name"].(string))
		if nil != err {
			logging.LOG.Errorf(
				"Cannot create etcd config file for node %v",
				map_ctx["etcd_name"])
			return false
		}
		ctx.Execute(writer, map_ctx)
		if !utils.SCPFiles([]string{writer.Name()},
			"/etc/etcd/etcd.conf", "file", true, conf.KUBERNETES_K8S_NODES[name]) {
			return false
		}
		ips = append(ips, conf.KUBERNETES_K8S_NODES[name])
	}

	source_ca_path := conf.CA_OUTPUT
	if !utils.SCPFiles([]string{source_ca_path},
		conf.ETCD_SSL, "", true, ips...) {
		return false
	}
	alias := fmt.Sprintf("alias etcdctl='etcdctl --endpoints=%s "+
		"--ca-file=%s/ca.pem --cert-file=%s/kubernetes.pem "+
		" --key-file=%s/kubernetes-key.pem'",
		etcd_endpoints, conf.ETCD_SSL, conf.ETCD_SSL, conf.ETCD_SSL)
	alias_cmd := `echo "` + alias + `" >> /root/.bashrc`
	cmd := "chown -R etcd:etcd " + conf.ETCD_SSL + ";" + alias_cmd +
		";systemctl enable etcd;systemctl restart etcd"

	return utils.RemoteCmd(cmd, ips...)
}
