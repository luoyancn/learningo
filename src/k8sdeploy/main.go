package k8sdeploy

import (
	_ "k8sdeploy/conf"
	"k8sdeploy/deploy"
	"k8sdeploy/logging"
	"k8sdeploy/utils"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var once sync.Once
var configfile string
var k8snodes []string
var k8snodeips []string

var rootcmd = &cobra.Command{
	Short: "Tools for deploy kubernetes clusters",
	Long:  ` The commands aims to deploy kubernetes clusters in easy way`,
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy the kubernetes cluster env",
	Long: `
Deploy the kubernetes cluster with golang tools.
`,
	PreRun: preparenv,
	Run:    deployk8s,
}

func init() {
	once.Do(func() {
		rootcmd.AddCommand(deployCmd)
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
	logging.GetLogger()
}

func preparenv(cmd *cobra.Command, args []string) {
	read_config()
	if 0 != os.Geteuid() {
		logging.LOG.Criticalf("Please execute this file with root permision\n")
		os.Exit(-1)
	}
	if !viper.IsSet("k8s.nodes") {
		msg := "Please ensure k8snodes confired in config files"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if !viper.IsSet("k8s.binary_path") {
		msg := "Please tell me where were your k8s binarys" +
			" in [k8s] section with binary_path\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}
	if !viper.IsSet("cfs.templates") {
		msg := "Please tell me where were your templates for generate ca" +
			" files in [cfs] section with templates\n"
		logging.LOG.Criticalf(msg)
		os.Exit(-1)
	}

	//k8snodes = viper.GetStringSlice("k8s.nodes")
	k8snode_map_ip := viper.GetStringMapString("k8s.nodes")
	for node, ip := range k8snode_map_ip {
		k8snodes = append(k8snodes, node)
		k8snodeips = append(k8snodeips, ip)
	}
	/*
		hostname, err := os.Hostname()
		if nil != err {
			logging.LOG.Critical(
				"Cannot get the hostname of this nodes:%v\n", err)
			os.Exit(-1)
		}
			logging.LOG.Noticef("%s\n", hostname)
			for _, host := range k8snodes {
				if hostname == host {
					break
				}
				logging.LOG.Noticef(
					"Please execute this programs on one of k8snodes\n")
				os.Exit(-1)
			}
	*/
	if !utils.SSHCheck(k8snodeips...) {
		logging.LOG.Criticalf(
			"Please ensure noauth-ssh configurated on all k8snodes\n")
		os.Exit(-1)
	}
}

func deployk8s(cmd *cobra.Command, args []string) {
	if !deploy.PrepareK8SBinary(k8snodeips...) {
		logging.LOG.Criticalf(
			"Failed to prepare k8s binary files on all k8snodes\n")
		os.Exit(-1)
	}

	if err := deploy.CreateCA(); nil != err {
		logging.LOG.Criticalf(
			"Failed to create CA files for k8snodes:%v\n", err)
		os.Exit(-1)
	}

	if !deploy.PrepareCAKey(k8snodeips...) {
		logging.LOG.Critical(
			"Failed to prepare ca-key files on all k8snodes\n")
		os.Exit(-1)
	}

	if !deploy.GenerateK8sCtx(k8snodeips...) {
		logging.LOG.Critical(
			"Failed to generate the kubernetes context on all k8snodes\n")
		os.Exit(-1)
	}

	if !deploy.GenerateK8sConfig(k8snodeips...) {
		logging.LOG.Critical(
			"Failed to generate the kubernetes config file on all k8snodes\n")
		os.Exit(-1)
	}
}

func Execute() {
	if err := rootcmd.Execute(); nil != err {
		os.Exit(1)
	}
}
