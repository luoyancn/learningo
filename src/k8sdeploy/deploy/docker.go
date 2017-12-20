package deploy

import (
	"k8sdeploy/conf"
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"path/filepath"
	"text/template"
)

func DeployDocker() bool {
	k8snodeips := []string{}
	for _, node := range conf.KUBERNETES_K8S_NODES {
		k8snodeips = append(k8snodeips, node)
	}
	funcMap := template.FuncMap{
		"minus": func(a, b int) int {
			return a - b
		},
	}
	_, template_name := filepath.Split(conf.DOCKER_TEMPLATE)
	ctx := template.Must(template.New(template_name).Funcs(
		funcMap).ParseFiles(conf.DOCKER_TEMPLATE))
	map_ctx := map[string]interface{}{
		"docker_hub_mirror":  conf.DOCKER_HUB_MIRROR,
		"insecure_registrys": conf.DOCKER_INSECURE_REGISTRYS}
	writer, err := os.Create("/tmp/docker.service")
	if nil != err {
		logging.LOG.Errorf(
			"Cannot create docker service config file:%v\n", err)
		return false
	}
	if err = ctx.Execute(writer, map_ctx); nil != err {
		logging.LOG.Errorf(
			"Cannot parse docker service config file:%v\n", err)
		return false
	}
	cmd := "systemctl daemon-reload;systemctl enable docker;systemctl restart docker"
	if utils.SCPFiles([]string{"/tmp/docker.service"},
		"/usr/lib/systemd/system/docker.service", "file", true, k8snodeips...) {
		return utils.RemoteCmd(cmd, k8snodeips...)
	}
	return false
}
