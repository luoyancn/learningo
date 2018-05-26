package rpc

import (
	"context"
	"fmt"
	"oceanstack/conf"
	"oceanstack/logging"
	"time"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
)

var Prefix = "etcd3_naming"
var deregister = make(chan struct{})

func register() error {
	serviceKey := fmt.Sprintf(
		"/%s/%s/%s", Prefix, conf.GRPC_ETCD_SERVICE_NAME,
		fmt.Sprintf("%s:%d", conf.GRPC_SERVER, conf.GRPC_PORT))

	logging.LOG.Infof("Service key is %s\n", serviceKey)
	tlsInfo := transport.TLSInfo{
		CertFile:      conf.GRPC_ETCD_CERT,
		KeyFile:       conf.GRPC_ETCD_KEY,
		TrustedCAFile: conf.GRPC_ETCD_CA,
	}

	logging.LOG.Infof("Tls info is %v\n", tlsInfo)
	_, err := tlsInfo.ClientConfig()
	if nil != err {
		logging.LOG.Errorf("Tls info is %v\n", tlsInfo)
		return err
	}

	logging.LOG.Infof("Endpoint is %v\n", conf.GRPC_ETCD_ENDPOINTS)
	client, err := etcd3.New(etcd3.Config{
		Endpoints: conf.GRPC_ETCD_ENDPOINTS,
		//TLS:       tlsConfig,
	})

	if err != nil {
		return fmt.Errorf("grpclb: create etcd3 client failed: %v", err)
	}

	//defer client.Close()
	resp, err := client.Grant(
		context.TODO(), int64(conf.GRPC_ETCD_TIMEOUT/time.Second+5))

	if err != nil {
		return fmt.Errorf("grpclb: create etcd3 lease failed: %v", err)
	}

	if _, err := client.Put(context.TODO(), serviceKey,
		fmt.Sprintf("%s:%d", conf.GRPC_SERVER, conf.GRPC_PORT),
		etcd3.WithLease(resp.ID)); err != nil {
		return fmt.Errorf(
			"grpclb: set service '%s' with ttl to etcd3 failed: %s",
			conf.GRPC_ETCD_SERVICE_NAME, err.Error())
	}

	if _, err := client.KeepAlive(context.TODO(), resp.ID); err != nil {
		return fmt.Errorf(
			"grpclb: refresh service '%s' with ttl to etcd3 failed: %s",
			conf.GRPC_ETCD_SERVICE_NAME, err.Error())
	}

	// wait deregister then delete
	go func() {
		<-deregister
		client.Delete(context.Background(), serviceKey)
		deregister <- struct{}{}
	}()

	return nil
}

func unRegister() {
	deregister <- struct{}{}
	<-deregister
}
