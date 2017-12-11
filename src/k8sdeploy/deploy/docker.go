package deploy

import (
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

func DeployDocker(k8snodeips ...string) bool {
	funcMap := template.FuncMap{
		"minus": func(a, b int) int {
			return a - b
		},
	}
	ctx := template.Must(template.New("docker.service.template").Funcs(
		funcMap).ParseFiles(viper.GetString("docker.template")))
	insecure_registrys := viper.GetStringSlice("docker.insecure_registrys")
	map_ctx := map[string]interface{}{
		"docker_hub_mirror":  viper.GetString("docker.docker_hub_mirror"),
		"insecure_registrys": insecure_registrys}
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
	ssh_key := utils.GenerateSshAuthConfig()
	result := make(chan bool, len(k8snodeips))
	if utils.SCPFiles([]string{"/tmp/docker.service"},
		"/usr/lib/systemd/system/docker.service", "", true, k8snodeips...) {
		for _, ip := range k8snodeips {
			go func(ip string) {
				ssh_conn, err := utils.GetSshConnection(ip, ssh_key)
				if nil != err {
					logging.LOG.Errorf(
						"Cannot connect to host %s:%v\n", ip, err)
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
						"Fail to execute command %s on host %s  :%v\n",
						cmd, ip, err)
					result <- false
					return
				}
				logging.LOG.Infof(
					"Sucess to execute command %s on host %s\n", cmd, ip)
				result <- true
			}(ip)
		}

		for i := 0; i < len(k8snodeips); i++ {
			if !<-result {
				return false
			}
		}
		return true
	}
	return false
}
