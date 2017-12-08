package deploy

import (
	"fmt"
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

func DeployEtcd(k8snodes map[string]string) bool {
	ctx := template.Must(template.ParseFiles(
		viper.GetString("etcd.template")))

	etcd_ssl := viper.GetString("etcd.ssl")
	map_ctx := map[string]interface{}{
		"etcd_protocal":    viper.GetString("etcd.protocal"),
		"etcd_token":       viper.GetString("etcd.cluster_token"),
		"client_cert_auth": viper.GetBool("etcd.client_cert_auth"),
		"peer_cert_auth":   viper.GetBool("etcd.peer_cert_auth"),
		"etcd_ssl":         etcd_ssl,
		"etcd_debug":       viper.GetBool("etcd.debug")}

	etcd_nodes := viper.GetStringMapString("etcd.nodes")
	nodes := []string{}
	ips := []string{}
	endpoints := []string{}
	for name, ip := range etcd_nodes {
		nodes = append(
			nodes, name+"="+map_ctx["etcd_protocal"].(string)+"://"+ip+":2380")
		endpoints = append(
			endpoints, map_ctx["etcd_protocal"].(string)+"://"+ip+":2379")
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
			"/etc/etcd/etcd.conf", "", true, k8snodes[name]) {
			return false
		}
		ips = append(ips, k8snodes[name])
	}

	source_ca_path := viper.GetString("cfs.output")
	if !utils.SCPFiles([]string{source_ca_path},
		viper.GetString("etcd.ssl"), "", true, ips...) {
		return false
	}
	ssh_key := utils.GenerateSshAuthConfig()
	result := make(chan bool, len(ips))
	alias := fmt.Sprintf("alias etcdctl='etcdctl --endpoints=%s "+
		"--ca-file=%s/ca.pem --cert-file=%s/kubernetes.pem "+
		" --key-file=%s/kubernetes-key.pem'",
		etcd_endpoints, etcd_ssl, etcd_ssl, etcd_ssl)
	alias_cmd := `echo "` + alias + `" >> /root/.bashrc`
	cmd := "chown -R etcd:etcd " + viper.GetString(
		"etcd.ssl") + ";" + alias_cmd +
		";systemctl enable ectd;systemctl restart etcd"
	for _, ip := range ips {
		go func(ip string) {
			ssh_conn, err := utils.GetSshConnection(ip, ssh_key)
			if nil != err {
				logging.LOG.Errorf("Cannot connect to host %s:%v\n", ip, err)
				result <- false
				return
			}
			defer ssh_conn.Close()
			session, err := ssh_conn.NewSession()
			if nil != err {
				logging.LOG.Errorf(
					"Cannot connect to host %s to exec:%v\n", ip, err)
				result <- false
				return
			}
			defer session.Close()
			logging.LOG.Noticef("Waiting to execute command:%s\n", cmd)
			if err = session.Run(cmd); nil != err {
				logging.LOG.Errorf(
					"Fail to change the file owner on host %s  :%v\n",
					ip, err)
				result <- false
				return
			}
			result <- true
		}(ip)
	}

	for i := 0; i < len(ips); i++ {
		if !<-result {
			return false
		}
	}
	return true
}
