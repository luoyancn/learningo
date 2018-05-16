package common

import (
	"oceanstack/logging"
	"oceanstack/rpc"
	"os"
	"os/signal"
	"syscall"
)

func Stop() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case s := <-sig:
		logging.LOG.Infof("Exit Ocean engine: Recived signal %s", s)
		if nil != rpc.GRPC {
			rpc.StopServer()
		}
		os.Exit(0)
	}
}
