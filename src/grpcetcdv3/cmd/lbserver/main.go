package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	rpc "grpcetcdv3/rpc/rpcserver"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startcmd = &cobra.Command{
	Use:   "start",
	Short: "Start lb server",
	Long:  ` Start lb server`,
	Run:   start,
	Args:  cobra.NoArgs,
}

var port int
var conf_path string
var once sync.Once

func init() {
	once.Do(func() {
		startcmd.PersistentFlags().IntVarP(
			&port, "port", "p", 8080,
			"The port of server listened, Default is 8080")
		startcmd.PersistentFlags().StringVarP(
			&conf_path, "config-path", "c", "",
			"The config file used for grpc server")
		viper.BindPFlag(
			"default.port", startcmd.PersistentFlags().Lookup("port"))
	})
}

func start(cmd *cobra.Command, args []string) {
	if "" != conf_path {
		viper.SetConfigFile(conf_path)
		if err := viper.ReadInConfig(); nil != err {
			fmt.Printf("Failed to read config file %s: %v\n", conf_path, err)
			os.Exit(-1)
		}
	}
	fmt.Printf("%d\n", viper.GetInt("default.port"))
	go rpc.StartServer(viper.GetInt("default.port"))
	wait()
}

func wait() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case s := <-sig:
		fmt.Printf("Terminating the lbserver: Recevied signal %v\n", s)
		rpc.StopServer()
		os.Exit(-1)
	}
}

func execute() {
	if err := startcmd.Execute(); nil != err {
		os.Exit(1)
	}
}

func main() {
	execute()
}
