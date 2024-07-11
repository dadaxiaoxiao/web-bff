package executor

import (
	"context"
	"fmt"
	"github.com/dadaxiaoxiao/web-bff/internal/domain"
)

// LocalFuncExecutor 本地执行器
type LocalFuncExecutor struct {
	// funcs 执行方法
	funcs map[string]func(ctx context.Context, job domain.Job) error
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}

// Exec 执行job
func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.Job) error {
	fn, ok := l.funcs[j.Name]
	if !ok {
		return fmt.Errorf("未知任务，你是否注册了？ %s", j.Name)
	}
	return fn(ctx, j)
}

// RegisterFunc 注册执行方法
// name 方法名
// fn 方法名对应的方法
func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, job domain.Job) error) {
	l.funcs[name] = fn
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: make(map[string]func(ctx context.Context, job domain.Job) error),
	}
}
