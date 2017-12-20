package k8sdeploy

import (
	"k8sdeploy/conf"
	"k8sdeploy/deploy"
	"k8sdeploy/logging"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var once sync.Once
var configfile string

var rootcmd = &cobra.Command{
	Short: "Tools for deploy kubernetes clusters",
	Long:  ` The commands aims to deploy kubernetes clusters in easy way`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initial the kubernetes cluster env",
	Long: `
Inital the kubernetes cluster with golang tools.
`,
	PreRun: preparenv,
	Run:    initk8senv,
}

var etcdCmd = &cobra.Command{
	Use:   "deploy-etcd",
	Short: "Deploy the etcd cluster",
	Long: `
Deploy the etcd cluster with golang tools.
`,
	PreRun: preparenv,
	Run:    initetcd,
}

var dockerCmd = &cobra.Command{
	Use:   "deploy-docker",
	Short: "Deploy the docker service",
	Long: `
Deploy the docker service with golang tools.
`,
	PreRun: preparenv,
	Run:    initdocker,
}

var k8smasterCmd = &cobra.Command{
	Use:   "deploy-k8s-master",
	Short: "Deploy the k8s master",
	Long: `
Deploy the k8s manster service with golang tools.
`,
	PreRun: preparenv,
	Run:    initk8smaster,
}

var k8snodeCmd = &cobra.Command{
	Use:   "deploy-k8s-node",
	Short: "Deploy the k8s ndoes",
	Long: `
Deploy the k8s nodes service with golang tools.
`,
	PreRun: preparenv,
	Run:    initk8snode,
}

var calicoCmd = &cobra.Command{
	Use:   "deploy-calico",
	Short: "Deploy the calico",
	Long: `
Deploy the calic network service with golang tools.
`,
	PreRun: preparenv,
	Run:    initcalico,
}

func init() {
	once.Do(func() {
		rootcmd.AddCommand(initCmd)
		rootcmd.AddCommand(etcdCmd)
		rootcmd.AddCommand(dockerCmd)
		rootcmd.AddCommand(k8smasterCmd)
		rootcmd.AddCommand(k8snodeCmd)
		rootcmd.PersistentFlags().StringVarP(
			&configfile, "config-file", "c", "",
			"The full path of config file")
		rootcmd.MarkPersistentFlagRequired("config-file")
	})
}

func read_config() {
	viper.SetConfigFile(configfile)
	if err := viper.ReadInConfig(); nil != err {
		log.Printf("ERROR:%v\n", err)
		os.Exit(-1)
	}
	conf.OverWriteConf()
	logging.GetLogger()

	if conf.DEBUG {
		for key, value := range viper.AllSettings() {
			settings := value.(map[string]interface{})
			for setting_key, setting_value := range settings {
				logging.LOG.Noticef(
					"%s.%s\t%v\n", key, setting_key, setting_value)
			}
		}
	}
}

func preparenv(cmd *cobra.Command, args []string) {
	read_config()
	if 0 != os.Geteuid() {
		logging.LOG.Criticalf("Please execute this file with root permision\n")
		os.Exit(-1)
	}
	if 0 >= len(conf.KUBERNETES_K8S_NODES) {
		msg := "Please ensure k8snodes confired in config files"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if "" == conf.KUBERNETES_K8S_BINARY {
		msg := "Please tell me where were your k8s binarys" +
			" in [kubernetes] section with k8s_binary\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if "" == conf.KUBERNETES_K8S_API_SERVER {
		msg := "Please tell me where were your k8s binarys" +
			" in [kubernetes] section with k8s_binary\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if "" == conf.KUBERNETES_K8S_API_SERVER_TEMPLATE {
		msg := "Please tell me where were your api server template" +
			" in [kubernetes] section with k8s_api_server_template\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if "" == conf.KUBERNETES_K8S_CONTROLLER_TEMPLATE {
		msg := "Please tell me where were your controller template file" +
			" in [kubernetes] section with k8s_controller_template\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if "" == conf.KUBERNETES_K8S_SCHEDULER_TEMPLATE {
		msg := "Please tell me where were your schduler template file" +
			" in [kubernetes] section with k8s_scheduler_template\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if "" == conf.KUBERNETES_KUBELET_TEMPLATE {
		msg := "Please tell me where were your kubelet template file" +
			" in [kubernetes] section with kubelet_template\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}

	// Check the required configurations in [ca] section
	if "" == conf.CA_TEMPLATE_PATH {
		msg := "Please tell me where were your templates for generate ca" +
			" files in [ca] section with template_path\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if _, err := os.Stat(conf.CA_TEMPLATE_PATH); os.IsNotExist(err) {
		logging.LOG.Criticalf(
			"Please ensure your ca template path %s exist\n",
			conf.CA_TEMPLATE_PATH)
		os.Exit(-1)
	}
	if "" == conf.CA_OUTPUT {
		msg := "Please tell me where to put your generated ca files" +
			" in [ca] section with output\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if _, err := os.Stat(conf.CA_OUTPUT); os.IsNotExist(err) {
		logging.LOG.Warningf(
			"The ca out put path %s not exist, try to creat it \n",
			conf.CA_TEMPLATE_PATH)
		err := os.MkdirAll(conf.CA_OUTPUT, 0700)
		if nil != err {
			logging.LOG.Criticalf("Cannot to create the ca out path:%v\n", err)
			os.Exit(-1)
		}
	}

	// Check the required configurations in [etcd] section
	if "" == conf.ETCD_TEMPLATE {
		msg := "Please tell me where your etcd template file" +
			" in [etcd] section with template\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if 0 == len(conf.ETCD_NODES) {
		msg := "Please tell me your etcd servers" +
			" in [etcd] section with nodes\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}

	// Check the required configurations in [docker] section
	if "" == conf.DOCKER_TEMPLATE {
		msg := "Please tell me where your docker template file" +
			" in [docker] section with template\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}

	// Check the required configurations in [calico] section
	if 0 == len(conf.CALICO_CNI_BINARY) {
		msg := "Please tell me your calico binary files" +
			" in [calico] section with cni_binary\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
}

func initk8senv(cmd *cobra.Command, args []string) {
	if !deploy.PrepareK8SBinary() {
		logging.LOG.Criticalf(
			"Failed to prepare k8s binary files on all k8snodes\n")
		os.Exit(-1)
	}
	if err := deploy.CreateCA(); nil != err {
		logging.LOG.Criticalf(
			"Failed to create CA files for k8snodes:%v\n", err)
		os.Exit(-1)
	}
	if !deploy.PrepareCAKey() {
		logging.LOG.Critical(
			"Failed to prepare ca-key files on all k8snodes\n")
		os.Exit(-1)
	}
	if !deploy.GenerateK8sCtx() {
		logging.LOG.Critical(
			"Failed to generate the kubernetes context on all k8snodes\n")
		os.Exit(-1)
	}

	if !deploy.GenerateK8sConfig() {
		logging.LOG.Critical(
			"Failed to generate the kubernetes config file on all k8snodes\n")
		os.Exit(-1)
	}
}

func initetcd(cmd *cobra.Command, args []string) {
	if !deploy.DeployEtcd() {
		logging.LOG.Critical(
			"Failed to deploy or init etcd cluster on all k8snodes\n")
		os.Exit(-1)
	}
}

func initdocker(cmd *cobra.Command, args []string) {
	if !deploy.DeployDocker() {
		logging.LOG.Critical(
			"Failed to generate the docker service config file on all k8snodes\n")
		os.Exit(-1)
	}
}

func initk8smaster(cmd *cobra.Command, args []string) {
	deploy.Deployk8sMaster()
}

func initk8snode(cmd *cobra.Command, args []string) {
	deploy.DeployK8sNode()
}

func initcalico(cmd *cobra.Command, args []string) {
}

func Execute() {
	if err := rootcmd.Execute(); nil != err {
		os.Exit(1)
	}
}
