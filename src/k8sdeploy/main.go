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
		logging.TRACE.Printf("Please execute this file with root permision\n")
		os.Exit(-1)
	}
	if !viper.IsSet("k8s.nodes") {
		msg := "Please ensure k8snodes confired in config files"
		logging.TRACE.Printf("%s\n", msg)
		os.Exit(-1)
	}
	if !viper.IsSet("k8s.binary_path") {
		msg := ` Please tell me where were your k8s binarys in
[k8s] section with binary_path\n`
		logging.TRACE.Printf("%s\n", msg)
		os.Exit(-1)
	}
	if !viper.IsSet("cfs.binary_path") {
		msg := ` Please tell me where were your cfs binarys in
[cfs] section with binary_path\n`
		logging.TRACE.Printf("%s\n", msg)
		os.Exit(-1)
	}
	if !viper.IsSet("cfs.templates") {
		msg := ` Please tell me where were your templates for generate ca files
in [cfs] section with templates\n`
		logging.TRACE.Printf("%s\n", msg)
		os.Exit(-1)
	}
	k8snodes = viper.GetStringSlice("k8s.nodes")
	if !utils.SSHCheck(k8snodes...) {
		logging.TRACE.Printf(
			"Please ensure noauth-ssh configurated on all k8snodes\n")
		os.Exit(-1)
	}
}

func deployk8s(cmd *cobra.Command, args []string) {
	if !deploy.PrepareK8SBinary(k8snodes...) {
		logging.TRACE.Printf(
			"Failed to prepare k8s binary files on all k8snodes\n")
		os.Exit(-1)
	}

	if err := deploy.CreateCA(); nil != err {
		logging.TRACE.Printf(
			"Failed to create CA files for k8snodes\n")
		os.Exit(-1)
	}
}

func Execute() {
	if err := rootcmd.Execute(); nil != err {
		os.Exit(1)
	}
}