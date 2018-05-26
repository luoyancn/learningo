package rpc

import (
	"errors"
	"fmt"
	"oceanstack/conf"
	"oceanstack/logging"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"google.golang.org/grpc/naming"
)

type resolver struct {
	serviceName string
}

func newResolver(serviceName string) *resolver {
	return &resolver{serviceName: serviceName}
}

func (this *resolver) Resolve(target string) (naming.Watcher, error) {

	logging.LOG.Infof("The target is :%s\n", target)
	if this.serviceName == "" {
		return nil, errors.New("grpclb: no service name provided")
	}

	tlsInfo := transport.TLSInfo{
		CertFile:      conf.GRPC_ETCD_CERT,
		KeyFile:       conf.GRPC_ETCD_KEY,
		TrustedCAFile: conf.GRPC_ETCD_CA,
	}

	_, err := tlsInfo.ClientConfig()
	if nil != err {
		return nil, err
	}

	client, err := etcd3.New(etcd3.Config{
		Endpoints: conf.GRPC_ETCD_ENDPOINTS,
		//TLS:       tlsConfig,
	})

	if nil != err {
		return nil, fmt.Errorf(
			"grpclb: creat etcd3 client failed: %s", err.Error())
	}
	return &watcher{resolver: this, client: *client}, nil
}
