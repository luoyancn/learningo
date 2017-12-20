package deploy

import (
	"k8sdeploy/conf"
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func Deployk8sMaster() bool {
	ips := []string{}
	for _, ip := range conf.KUBERNETES_K8S_NODES {
		ips = append(ips, ip)
	}
	k8sapis := conf.ETCD_NODES
	etcd_servers_array := []string{}
	for _, ip := range k8sapis {
		etcd_server := conf.ETCD_PROTOCAL + "://" + ip + ":2379"
		etcd_servers_array = append(etcd_servers_array, etcd_server)
	}
	etcd_servers := strings.Join(etcd_servers_array, ",")
	api_ctx := template.Must(
		template.ParseFiles(conf.KUBERNETES_K8S_API_SERVER_TEMPLATE))
	controller_ctx := template.Must(
		template.ParseFiles(conf.KUBERNETES_K8S_CONTROLLER_TEMPLATE))
	scheduler_ctx := template.Must(
		template.ParseFiles(conf.KUBERNETES_K8S_SCHEDULER_TEMPLATE))

	cluster_service_ip_cidr := conf.KUBERNETES_K8S_CLUSTER_SERVICE_IP_CIDR

	api_map := map[string]interface{}{
		"etcd_servers":             etcd_servers,
		"apiserver_count":          len(k8sapis),
		"apiserver_insecure_port":  conf.KUBERNETES_K8S_APISERVER_INSECURE_PORT,
		"apiserver_runtime_config": conf.KUBERNETES_K8S_APISERVER_RUNTIME_CONFIG,
		"apiserver_secure_port":    conf.KUBERNETES_K8S_APISERVER_SECURE_PORT,
		"cluster_service_ip_cidr":  cluster_service_ip_cidr,
		"service_node_port_range":  conf.KUBERNETES_K8S_SERVICE_NODE_PORT_RANGE}

	insecure_apiserver := "http://" + conf.KUBERNETES_K8S_API_SERVER +
		":" + strconv.Itoa(conf.KUBERNETES_K8S_APISERVER_INSECURE_PORT)
	controller_map := map[string]interface{}{
		"cluster_pod_ip_cidr":     conf.KUBERNETES_K8S_CLUSTER_POD_IP_CIDR,
		"cluster_service_ip_cidr": cluster_service_ip_cidr,
		"insecure_apiserver":      insecure_apiserver,
		"controller_manager_port": conf.KUBERNETES_K8S_CONTROLLER_MANAGER_PORT,
		"cluster_name":            conf.KUBERNETES_K8S_CLUSTER_NAME}

	scheduler_map := map[string]interface{}{
		"scheduler_port":     conf.KUBERNETES_K8S_SCHEDULER_PORT,
		"insecure_apiserver": insecure_apiserver}

	controller_writer, err := os.Create(
		"/tmp/kube-controller-manager.service")
	if nil != err {
		logging.LOG.Errorf(
			"Cannot create kube controller config file:%v\n", err)
		return false
	}
	if err = controller_ctx.Execute(
		controller_writer, controller_map); nil != err {
		logging.LOG.Errorf(
			"Cannot parse kube controller config file:%v\n", err)
		return false
	}

	scheduler_writer, err := os.Create(
		"/tmp/kube-scheduler.service")
	if nil != err {
		logging.LOG.Errorf(
			"Cannot create kube scheduler config file:%v\n", err)
		return false
	}
	if err = scheduler_ctx.Execute(
		scheduler_writer, scheduler_map); nil != err {
		logging.LOG.Errorf(
			"Cannot parse kube scheduler config file:%v\n", err)
		return false
	}

	for node, api_server := range k8sapis {
		api_map["api_server"] = api_server
		api_writer, err := os.Create("/tmp/kube-apiserver.service." + node)
		if nil != err {
			logging.LOG.Errorf(
				"Cannot create kube apiserver config file:%v\n", err)
			return false
		}
		if err = api_ctx.Execute(api_writer, api_map); nil != err {
			logging.LOG.Errorf(
				"Cannot parse kube apiserver config file:%v\n", err)
			return false
		}
		if !utils.SCPFiles([]string{api_writer.Name()},
			"/usr/lib/systemd/system/kube-apiserver.service",
			"file", true, conf.KUBERNETES_K8S_NODES[node]) {
			return false
		}
	}
	if !utils.SCPFiles([]string{controller_writer.Name(),
		scheduler_writer.Name()}, "/usr/lib/systemd/system/",
		"", true, ips...) {
		return false
	}
	cmd := "mkdir -p /var/log/kubernetes;for id in kube-{apiserver,controller-manager,scheduler};do systemctl daemon-reload;systemctl enable $id;systemctl restart $id;done"
	return utils.RemoteCmd(cmd, ips...)
}
