package ioc

import (
	cronjobv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/cronjob/v1"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/web-bff/internal/job"
	"github.com/dadaxiaoxiao/web-bff/internal/job/executor"
)

func InitHttpExecutor() *executor.HttpExecutor {
	return executor.NewHttpExecutor()
}

// InitLocalFuncExecutor 初始化 本地执行器
func InitLocalFuncExecutor() *executor.LocalFuncExecutor {
	res := executor.NewLocalFuncExecutor()
	return res
}

// InitScheduler 初始化 任务调度
func InitScheduler(jobClient cronjobv1.CronJobServiceClient,
	local *executor.LocalFuncExecutor,
	http *executor.HttpExecutor,
	l accesslog.Logger) *job.Scheduler {
	res := job.NewScheduler(jobClient, l)
	// 注册 Executor
	res.RegisterExecutor(local)
	res.RegisterExecutor(http)
	return res
}
