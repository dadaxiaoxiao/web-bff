package executor

import (
	"context"
	"github.com/dadaxiaoxiao/web-bff/internal/domain"
)

// Executor 执行器
type Executor interface {

	// Name 执行器名称
	Name() string

	// Exec 执行 Job
	// ctx 是整个任务调度的上下文
	// 具体实现来控制
	Exec(ctx context.Context, j domain.Job) error
}

