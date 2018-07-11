package etcdv3

import (
	"fmt"
	"strings"

	etcd "github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
)

var deregister = make(chan struct{})

type Option struct {
	RegistryDir string
	ServiceName string
	NodeID      string
	NData       string
}

func Register(config etcd.Config, opt Option) error {
	key := strings.Join(
		[]string{opt.RegistryDir, opt.ServiceName, opt.NodeID}, "/")
	client, err := etcd.New(config)
	if err != nil {
		panic(err)
	}
	go func() {
		<-deregister
		fmt.Printf("Delete %s from etcd\n", key)
		client.Delete(context.Background(), key)
		deregister <- struct{}{}
	}()

	resp, err := client.Grant(context.TODO(), int64(10))
	if err != nil {
		return fmt.Errorf("grpclb: create etcd3 lease failed: %v", err)
	}

	if _, err = client.Put(context.TODO(), key,
		opt.NData, etcd.WithLease(resp.ID)); err != nil {
		return fmt.Errorf(
			"grpclb: set service '%s' with ttl to etcd3 failed: %s",
			"local lb", err.Error())
	}

	if _, err := client.KeepAlive(context.TODO(), resp.ID); err != nil {
		return fmt.Errorf(
			"grpclb: refresh service '%s' with ttl to etcd3 failed: %s",
			"local lb", err.Error())
	}

	return nil
}

func UnRegister() {
	deregister <- struct{}{}
	<-deregister
}
