package deploy

import (
	"k8sdeploy/utils"

	"github.com/spf13/viper"
)

func PrepareK8SBinary(nodes ...string) bool {
	binary_files_path := viper.GetString("k8s.binary_path")
	dest_binary_path := viper.GetString("k8s.target_path")
	overwrite_k8s_binary := viper.GetBool("k8s.overwrite_binary")
	return utils.SCPFiles([]string{binary_files_path},
		dest_binary_path, "binary", overwrite_k8s_binary, nodes...)
}

func PrepareCAKey(nodes ...string) bool {
	source_json_path := viper.GetString("cfs.templates")
	source_ca_path := viper.GetString("cfs.output")
	dest_ca_path := viper.GetString("k8s.ssl_config_path")
	overwrite_k8s_ssl := viper.GetBool("k8s.overwrite_ssl")
	return utils.SCPFiles([]string{source_json_path, source_ca_path},
		dest_ca_path, "", overwrite_k8s_ssl, nodes...)
}
