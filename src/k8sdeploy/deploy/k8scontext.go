package deploy

import (
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

func decide_cmd_and_path(config_target_path string,
	config_source_path string, cmd_or_file_name string,
	cmd_or_file_path *string) {
	target_path := filepath.Join(
		viper.GetString(config_target_path), cmd_or_file_name)
	if _, err := os.Stat(target_path); os.IsNotExist(err) {
		*cmd_or_file_path = filepath.Join(
			viper.GetString(config_source_path), cmd_or_file_name)
	} else {
		*cmd_or_file_path = target_path
	}
}

func GenerateK8sCtx(k8snodes ...string) bool {
	var kubectl_cmd string
	decide_cmd_and_path("k8s.target_path",
		"k8s.binary_path", "kubectl", &kubectl_cmd)

	var ca_pem string
	decide_cmd_and_path("k8s.ssl_config_path",
		"cfs.output", "ca.pem", &ca_pem)

	var admin_ca_pem string
	decide_cmd_and_path("k8s.ssl_config_path",
		"cfs.output", "admin-ca.pem", &admin_ca_pem)

	var admin_key_pem string
	decide_cmd_and_path("k8s.ssl_config_path",
		"cfs.output", "admin-key.pem", &admin_key_pem)

	cluster_name := viper.GetString("k8s.cluster_name")
	set_cluster_cmd := exec.Command(kubectl_cmd, "config", "set-cluster",
		cluster_name, "--embed-certs=true",
		"--certificate-authority="+ca_pem,
		"--server="+viper.GetString("k8s.api_server"))
	logging.LOG.Infof("Running the command :%v\n", set_cluster_cmd.Args)
	if err := set_cluster_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to create the kubectl config file:%v\n", err)
		return false
	}

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

	set_context_cmd := exec.Command(
		kubectl_cmd, "config", "set-context", cluster_name,
		"--cluster="+cluster_name, "--user=admin")
	logging.LOG.Infof("Running the command :%v\n", set_context_cmd.Args)
	if err := set_context_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to set context to the kubectl config file:%v\n", err)
		return false
	}

	use_context_cmd := exec.Command(
		kubectl_cmd, "config", "use-context", cluster_name)
	logging.LOG.Infof("Running the command :%v\n", use_context_cmd.Args)
	if err := use_context_cmd.Start(); nil != err {
		logging.LOG.Fatalf(
			"Failed to switch context to the admin context:%v\n", err)
		return false
	}
	return utils.SCPFiles([]string{"/root/.kube"},
		"/root/.kube", "", true, k8snodes...)
}
