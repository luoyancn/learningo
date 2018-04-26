package main

import (
	"fmt"
	"oceanstack/common"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var once sync.Once
var configfile string

var rootcmd = &cobra.Command{
	Short: "Client of ocean stack",
	Long:  ` The commands aims to support client for ocean client`,
}

var envcmd = &cobra.Command{
	Use:   "env",
	Short: "Ocean env commands collections",
	Long:  ` Ocean env commands collections`,
}

var vercmd = &cobra.Command{
	Use:   "version",
	Short: "Get the oceanclient version",
	Long:  ` Get the version of oceanclient binary`,
	Run:   get_version,
	Args:  cobra.NoArgs,
}

var initenv = &cobra.Command{
	Use:   "init",
	Short: "Init the ocean evn",
	Long:  ` Init the ocean env`,
	Run:   init_env,
	Args:  cobra.NoArgs,
}

func init() {
	once.Do(func() {
		envcmd.PersistentFlags().StringVarP(
			&configfile, "config-file", "c", "",
			"The full path of config file (Required)")
		envcmd.MarkPersistentFlagRequired("config-file")
		envcmd.AddCommand(initenv)
		rootcmd.AddCommand(envcmd)
		rootcmd.AddCommand(vercmd)
	})
}

func init_env(cmd *cobra.Command, args []string) {
	fmt.Println("hello")
}

func get_version(cmd *cobra.Command, args []string) {
	common.Versions("OceanClient")
}

func Execute() {
	if err := rootcmd.Execute(); nil != err {
		os.Exit(1)
	}
}

func main() {
	Execute()
}
