package ioc

import (
	cronjobv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/cronjob/v1"
	"github.com/dadaxiaoxiao/go-pkg/grpcx/interceptors/trace"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitCronJobGRPCClient(etcdClient *etcdv3.Client, redisClient redis.Cmdable) cronjobv1.CronJobServiceClient {
	type Config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.cronjob", &cfg)
	if err != nil {
		panic(err)
	}

	// 下面是服务发现
	rs, err := resolver.NewBuilder(etcdClient)
	opts := []grpc.DialOption{
		grpc.WithResolvers(rs),
		// 拦截器
		grpc.WithChainUnaryInterceptor(
			// 链路追踪
			trace.NewInterceptorBuilder(nil, nil).BuildClient(),
		),
	}

	if cfg.Secure {
		// 加载证书之类的东西
		// 启动 Https
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.Dial(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	return cronjobv1.NewCronJobServiceClient(cc)
}
