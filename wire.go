//go:build wireinject

package main

import (
	"github.com/dadaxiaoxiao/web-bff/ioc"
	"github.com/google/wire"
)

var thirdPartyProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitRedis,
	ioc.InitOTEL,
	ioc.InitEtcdClient,
)

var cronjobProvider = wire.NewSet(
	ioc.InitCronJobGRPCClient,
	ioc.InitHttpExecutor,
	ioc.InitLocalFuncExecutor,
	ioc.InitScheduler,
)

func InitApp() *WebApp {
	wire.Build(
		thirdPartyProvider,
		cronjobProvider,
		wire.Struct(new(WebApp), "Scheduler"),
	)
	return new(WebApp)
}
