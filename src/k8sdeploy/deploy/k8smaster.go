package deploy

import (
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

func Deployk8sMaster(k8snodes map[string]string) bool {
	ips := []string{}
	for _, ip := range k8snodes {
		ips = append(ips, ip)
	}
	k8sapis := viper.GetStringMapString("etcd.nodes")
	etcd_servers_array := []string{}
	for _, ip := range k8sapis {
		etcd_server := "https://" + ip + ":2379"
		etcd_servers_array = append(etcd_servers_array, etcd_server)
	}
	etcd_servers := strings.Join(etcd_servers_array, ",")
	api_ctx := template.Must(
		template.ParseFiles(viper.GetString("k8s.api_server_template")))
	controller_ctx := template.Must(
		template.ParseFiles(viper.GetString("k8s.controller_template")))
	scheduler_ctx := template.Must(
		template.ParseFiles(viper.GetString("k8s.scheduler_template")))

	cluster_service_ip_cidr := viper.GetString("k8s.cluster_service_ip_cidr")

	api_map := map[string]interface{}{
		"etcd_servers":             etcd_servers,
		"apiserver_count":          len(k8sapis),
		"apiserver_insecure_port":  viper.GetString("k8s.apiserver_insecure_port"),
		"apiserver_runtime_config": viper.GetString("k8s.apiserver_runtime_config"),
		"apiserver_secure_port":    viper.GetString("k8s.apiserver_secure_port"),
		"cluster_service_ip_cidr":  cluster_service_ip_cidr,
		"service_node_port_range":  viper.GetString("k8s.service_node_port_range")}

	insecure_apiserver := "http://" + viper.GetString("k8s.api_server") +
		":" + viper.GetString("k8s.api_insecure_port")
	controller_map := map[string]interface{}{
		"cluster_pod_ip_cidr":     viper.GetString("k8s.cluster_pod_ip_cidr"),
		"cluster_service_ip_cidr": cluster_service_ip_cidr,
		"insecure_apiserver":      insecure_apiserver,
		"controller_manager_port": viper.GetInt("k8s.controller_manager_port"),
		"cluster_name":            viper.GetString("k8s.cluster_name")}

	scheduler_map := map[string]interface{}{
		"scheduler_port":     viper.GetString("k8s.scheduler_port"),
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
			"file", true, k8snodes[node]) {
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
