package deploy

import (
	"k8sdeploy/conf"
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

func GenerateK8sCtx() bool {
	k8snodes := []string{}
	for _, node := range conf.KUBERNETES_K8S_NODES {
		k8snodes = append(k8snodes, node)
	}
	kubectl_cmd := filepath.Join(conf.KUBERNETES_K8S_BINARY, "kubectl")
	ca_pem := filepath.Join(conf.CA_OUTPUT, "ca.pem")
	admin_ca_pem := filepath.Join(conf.CA_OUTPUT, "admin.pem")
	admin_key_pem := filepath.Join(conf.CA_OUTPUT, "admin-key.pem")

	cluster_name := conf.KUBERNETES_K8S_CLUSTER_NAME
	set_cluster_cmd := exec.Command(kubectl_cmd, "config", "set-cluster",
		cluster_name, "--embed-certs=true",
		"--certificate-authority="+ca_pem,
		"--server=https://"+conf.KUBERNETES_K8S_API_SERVER+":"+
			strconv.Itoa(conf.KUBERNETES_K8S_APISERVER_SECURE_PORT))
	logging.LOG.Infof("Running the command :%v\n", set_cluster_cmd.Args)
	if err := set_cluster_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to create the kubectl config file:%v\n", err)
		return false
	}

	time.Sleep(5 * time.Millisecond)
	set_credentials_cmd := exec.Command(
		kubectl_cmd, "config", "set-credentials", "admin",
		"--embed-certs=true", "--client-certificate="+admin_ca_pem,
		"--client-key="+admin_key_pem)
	logging.LOG.Infof("Running the command :%v\n", set_credentials_cmd.Args)
	if err := set_credentials_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to set credentials to the kubectl config file:%v\n", err)
		return false
	}

	time.Sleep(5 * time.Millisecond)
	set_context_cmd := exec.Command(
		kubectl_cmd, "config", "set-context", cluster_name,
		"--cluster="+cluster_name, "--user=admin")
	logging.LOG.Infof("Running the command :%v\n", set_context_cmd.Args)
	if err := set_context_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to set context to the kubectl config file:%v\n", err)
		return false
	}

	time.Sleep(5 * time.Millisecond)
	use_context_cmd := exec.Command(
		kubectl_cmd, "config", "use-context", cluster_name)
	logging.LOG.Infof("Running the command :%v\n", use_context_cmd.Args)
	if err := use_context_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to switch context to the admin context:%v\n", err)
		return false
	}
	create_channel := make(chan struct{})
	go func(f string) {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		for _ = range ticker.C {
			if _, err := os.Stat(f); nil != err {
				logging.LOG.Warningf(
					"Waiting the kubenetes file %s created end\n", f)
			} else {
				logging.LOG.Noticef("kubenetes file %s created end\n", f)
				break
			}
		}
		create_channel <- struct{}{}
	}("/root/.kube/config")

	<-create_channel

	return utils.SCPFiles([]string{"/root/.kube"},
		"/root/.kube", "", true, k8snodes...)
}
