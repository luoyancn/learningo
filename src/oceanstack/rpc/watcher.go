package rpc

import (
	"context"
	"fmt"
	"oceanstack/logging"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/naming"
)

type watcher struct {
	resolver  *resolver
	client    etcd3.Client
	is_inited bool
}

func (this *watcher) Close() {}

func (this *watcher) Next() ([]*naming.Update, error) {
	prefix := fmt.Sprintf("/%s/%s/", Prefix, this.resolver.serviceName)
	logging.LOG.Infof("Prefix is :%s\n", prefix)
	if !this.is_inited {
		resp, err := this.client.Get(
			context.Background(), prefix, etcd3.WithPrefix())
		logging.LOG.Infof("Response is %v\n", resp)
		this.is_inited = true
		if nil == err {
			addrs := extractAddrs(resp)
			if count := len(addrs); count != 0 {
				update := make([]*naming.Update, count)
				for i := range addrs {
					update[i] = &naming.Update{Op: naming.Add, Addr: addrs[i]}
				}
				return update, nil
			}
		}
	}

	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()
	rch := this.client.Watch(ctx, prefix, etcd3.WithPrefix())

	for etcd_resp := range rch {
		for _, ev := range etcd_resp.Events {
			switch ev.Type {
			case mvccpb.PUT:
				return []*naming.Update{{Op: naming.Add,
					Addr: string(ev.Kv.Value)}}, nil
			case mvccpb.DELETE:
				return []*naming.Update{{Op: naming.Delete,
					Addr: string(ev.Kv.Value)}}, nil
			}
		}
	}

	return nil, nil
}

func extractAddrs(resp *etcd3.GetResponse) []string {
	addrs := []string{}

	if resp == nil || resp.Kvs == nil {
		return addrs
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			addrs = append(addrs, string(v))
		}
	}

	return addrs
}
