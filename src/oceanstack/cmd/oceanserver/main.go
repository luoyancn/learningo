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
	Short: "Server of ocean stack",
	Long:  ` This is the binary of ocean server`,
}

var startcmd = &cobra.Command{
	Use:   "start",
	Short: "Start Ocean server",
	Long:  ` Start Ocean server`,
	Run:   serve,
}

var vercmd = &cobra.Command{
	Use:   "version",
	Short: "Get the ocean server version",
	Long:  ` Get the version of ocean server binary`,
	Run:   get_version,
	Args:  cobra.NoArgs,
}

func init() {
	once.Do(func() {
		startcmd.PersistentFlags().StringVarP(
			&configfile, "config-file", "c", "",
			"The full path of config file (Required)")
		startcmd.MarkPersistentFlagRequired("config-file")
		rootcmd.AddCommand(startcmd)
		rootcmd.AddCommand(vercmd)
	})
}

func serve(cmd *cobra.Command, args []string) {
	fmt.Println("hello")
}

func get_version(cmd *cobra.Command, args []string) {
	common.Versions("OceanServer")
}

func Execute() {
	if err := rootcmd.Execute(); nil != err {
		os.Exit(1)
	}
}

func main() {
	Execute()
}
