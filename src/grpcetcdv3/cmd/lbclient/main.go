package main

import (
	"fmt"
	"os"
	"sync"

	"grpcetcdv3/rpc"

	"github.com/spf13/cobra"
)

var startcmd = &cobra.Command{
	Use:   "start",
	Short: "Start lb client",
	Long:  ` Start lb client`,
	Run:   start,
	Args:  cobra.NoArgs,
}

var once sync.Once

func init() {
	once.Do(func() {
		rpc.InitGrpcClientPool()
	})
}

func start(cmd *cobra.Command, args []string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%s\n", rpc.Call())
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
