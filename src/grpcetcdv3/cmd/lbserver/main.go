package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"grpcetcdv3/rpc"

	"github.com/spf13/cobra"
)

var startcmd = &cobra.Command{
	Use:   "start",
	Short: "Start lb server",
	Long:  ` Start lb server`,
	Run:   start,
	Args:  cobra.NoArgs,
}

var port int
var once sync.Once

func init() {
	once.Do(func() {
		startcmd.PersistentFlags().IntVarP(
			&port, "port", "p", 8080,
			"The port of server listened, Default is 8080")
	})
}

func start(cmd *cobra.Command, args []string) {
	go rpc.StartServer(port)
	wait()
}

func wait() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-sig:
		rpc.StopServer()
		os.Exit(0)
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
