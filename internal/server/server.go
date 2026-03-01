package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer, NewRegistrar)

func NewRegistrar(conf *conf.Registry) registry.Registrar {
	endpoints := conf.Etcd.GetEndpoints()
	if endpoints == nil || len(endpoints[0]) == 0 {
		panic(errors.New("etcd endpoints is empty"))
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		panic(err)
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	// 创建 client 后立即测试连接
	if _, err := cli.Status(ctx, endpoints[0]); err != nil {
		panic(fmt.Errorf("etcd connection failed: %w", err))
	}

	r := etcd.New(cli)
	//c := consulAPI.DefaultConfig()
	//c.Address = conf.Consul.Address
	//c.Scheme = conf.Consul.Scheme
	//cli, err := consulAPI.NewClient(c)
	//if err != nil {
	//	panic(err)
	//}
	//r := consul.New(cli, consul.WithHealthCheck(false))
	return r
}
