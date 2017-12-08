package deploy

import (
	"text/template"

	"github.com/spf13/viper"
)

func DeployDocker(k8snodes map[string]string) bool {
	ctx := template.Must(template.ParseFiles(
		viper.GetString("docker.template")))
	map_ctx := map[string]string{
		"docker_hub_mirror":    viper.GetString("docker.docker_hub_mirror"),
		"docker_gcr_io_mirror": viper.GetString("docker.docker_gcr_io_mirror")}
	harbor_registry := viper.GetString("docker.harbor_registry")
	if "" != harbor_registry {
		map_ctx["harbor_registry"] = harbor_registry
	}
}
