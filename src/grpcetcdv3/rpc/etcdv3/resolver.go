package etcdv3

import (
	"errors"
	"fmt"

	etcd "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/naming"
)

type etcdResolver struct {
	config      etcd.Config
	registryDir string
	serviceName string
}

func NewResolver(registryDir string, serviceName string,
	cfg etcd.Config) naming.Resolver {
	return &etcdResolver{registryDir: registryDir,
		serviceName: serviceName, config: cfg}
}

func (this *etcdResolver) Resolve(target string) (naming.Watcher, error) {
	if this.serviceName == "" {
		return nil, errors.New("no service name provided")
	}
	client, err := etcd.New(this.config)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s/%s", this.registryDir, this.serviceName)
	return newEtcdWatcher(key, client), nil
}
