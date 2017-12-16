package deploy

import (
	"io/ioutil"
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func random_bootstrap_token() string {
	head_cmd := "head -c 16 /dev/urandom |od -An -t x |tr -d ' '"
	out, _ := exec.Command("bash", "-c", head_cmd).Output()
	return strings.Replace(string(out), "\n", "", 1)
}

func generate_token_csv(k8snodes ...string) (bool, string) {
	bootstrap_token := random_bootstrap_token()
	ctx := []string{bootstrap_token, "kubelet-bootstrap",
		"10001", `"system:kubelet-bootstrap"`}
	token_csv := viper.GetString("k8s.tmp") + "/token.csv"
	err := ioutil.WriteFile(token_csv, []byte(strings.Join(ctx, ",")), 0644)
	if nil != err {
		logging.LOG.Errorf("Cannot create the token csv file:%v\n", err)
		return false, ""
	}
	return utils.SCPFiles(
		[]string{token_csv}, viper.GetString("k8s.config_path"),
		"", true, k8snodes...), bootstrap_token
}

func GenerateK8sConfig(k8snodes ...string) bool {
	ok, bootstrap_token := generate_token_csv(k8snodes...)
	if !ok {
		return false
	}
	var kubectl_cmd string
	decide_cmd_and_path("k8s.target_path",
		"k8s.binary_path", "kubectl", &kubectl_cmd)

	var ca_pem string
	decide_cmd_and_path("k8s.ssl_config_path",
		"cfs.output", "ca.pem", &ca_pem)

	var kube_proxy_pem string
	decide_cmd_and_path("k8s.ssl_config_path",
		"cfs.output", "kube-proxy.pem", &kube_proxy_pem)

	var kube_proxy_key_pem string
	decide_cmd_and_path("k8s.ssl_config_path",
		"cfs.output", "kube-proxy-key.pem", &kube_proxy_key_pem)

	cluster_name := viper.GetString("k8s.cluster_name")

	kube_file := "/tmp/kubelet-bootstrap.kubeconfig"
	kubeconfig := "--kubeconfig=" + kube_file
	kube_proxy := "/tmp/kube-proxy.kubeconfig"
	kubeproxyconfig := "--kubeconfig=" + kube_proxy

	set_cluster_cmd := exec.Command(kubectl_cmd, "config", "set-cluster",
		cluster_name, "--embed-certs=true",
		"--server=https://"+viper.GetString("k8s.api_server")+":"+
			viper.GetString("k8s.apiserver_secure_port"),
		"--certificate-authority="+ca_pem, kubeconfig)
	logging.LOG.Infof("Running the command :%v\n", set_cluster_cmd.Args)
	if err := set_cluster_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to set kubelet cluster:%v\n", err)
		return false
	}

	set_credentials_cmd := exec.Command(
		kubectl_cmd, "config", "set-credentials",
		"--token="+bootstrap_token, kubeconfig)
	logging.LOG.Infof("Running the command :%v\n", set_credentials_cmd.Args)
	if err := set_credentials_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to set kubelet cluster credentials:%v\n", err)
		return false
	}

	set_context_cmd := exec.Command(
		kubectl_cmd, "config", "set-context",
		"--cluster="+cluster_name, "--user=kubelet-bootstrap", kubeconfig)
	logging.LOG.Infof("Running the command :%v\n", set_context_cmd.Args)
	if err := set_context_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to set kubelet cluster context:%v\n", err)
		return false
	}

	use_context_cmd := exec.Command(
		kubectl_cmd, "config", "use-context", "default",
		"--cluster="+cluster_name, "--user=kubelet-bootstrap", kubeconfig)
	logging.LOG.Infof("Running the command :%v\n", use_context_cmd.Args)
	if err := use_context_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to switch kubelet cluster context:%v\n", err)
		return false
	}

	set_proxy_cluster_cmd := exec.Command(
		kubectl_cmd, "config", "set-cluster", cluster_name,
		"--embed-certs=true", "--server=https://"+
			viper.GetString("k8s.api_server")+
			":"+viper.GetString("k8s.apiserver_secure_port"),
		"--certificate-authority="+ca_pem, kubeproxyconfig)
	logging.LOG.Infof("Running the command :%v\n", set_proxy_cluster_cmd.Args)
	if err := set_proxy_cluster_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to set kubeproxy cluster:%v\n", err)
		return false
	}

	set_proxy_cred_cmd := exec.Command(
		kubectl_cmd, "config", "set-credentials", "kube-proxy",
		"--embed-certs=true", "--client-certificate="+kube_proxy_pem,
		"--client-key="+kube_proxy_key_pem, kubeproxyconfig)
	logging.LOG.Infof("Running the command :%v\n", set_proxy_cred_cmd.Args)
	if err := set_proxy_cred_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to set kubeproxy cluster credentials:%v\n", err)
		return false
	}

	set_proxy_context_cmd := exec.Command(
		kubectl_cmd, "config", "set-context", "default",
		"--cluster="+cluster_name, "--user=kube-proxy", kubeproxyconfig)
	logging.LOG.Infof("Running the command :%v\n", set_proxy_context_cmd.Args)
	if err := set_proxy_context_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to set kubeproxy cluster context:%v\n", err)
		return false
	}

	use_proxy_context_cmd := exec.Command(
		kubectl_cmd, "config", "use-context", "default", kubeproxyconfig)
	logging.LOG.Infof("Running the command :%v\n", use_proxy_context_cmd.Args)
	if err := use_proxy_context_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to switch kubeproxy cluster context:%v\n", err)
		return false
	}

	runtime.GOMAXPROCS(2)
	create_channel := make(chan struct{}, 2)

	for _, f := range []string{kube_file, kube_proxy} {
		go func(f string) {
			ticker := time.NewTicker(time.Millisecond * 10)
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
		}(f)
	}

	for i := 0; i < 2; i++ {
		<-create_channel
	}

	return utils.SCPFiles(
		[]string{kube_file, kube_proxy},
		viper.GetString("k8s.config_path"), "", true, k8snodes...)
}
