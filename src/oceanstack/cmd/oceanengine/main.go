package main

import (
	"oceanstack/common"
	"oceanstack/logging"
	"oceanstack/rpc"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var once sync.Once
var configfile string

var rootcmd = &cobra.Command{
	Short: "Server of ocean stack engine",
	Long:  ` This is the binary of ocean engine`,
}

var startcmd = &cobra.Command{
	Use:   "start",
	Short: "Start Ocean engine",
	Long:  ` Start Ocean engine`,
	Run:   start,
	Args:  cobra.NoArgs,
}

var vercmd = &cobra.Command{
	Use:   "version",
	Short: "Get the ocean engine version",
	Long:  ` Get the version of ocean engine binary`,
	Run:   get_version,
	Args:  cobra.NoArgs,
}

func init() {
	once.Do(func() {
		startcmd.PersistentFlags().StringVarP(
			&configfile, "config-file", "c", "",
			"The full path of config file (Required, Only yaml, json and toml supported)")
		startcmd.MarkPersistentFlagRequired("config-file")
		rootcmd.AddCommand(startcmd)
		rootcmd.AddCommand(vercmd)
	})
}

func start(cmd *cobra.Command, args []string) {
	common.ReadConfig(configfile, "oceanengine", logging.FILE_ENABLED)
	go rpc.StartServer()
	logging.LOG.Infof("Ocean engine started\n")
	common.Wait()
}

func get_version(cmd *cobra.Command, args []string) {
	common.Versions("OceanEngine")
}

func execute() {
	if err := rootcmd.Execute(); nil != err {
		os.Exit(1)
	}
}

func main() {
	execute()
}
