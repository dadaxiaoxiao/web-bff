package job

import (
	"context"
	cronjobv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/cronjob/v1"
	"github.com/dadaxiaoxiao/go-pkg/accesslog"
	"github.com/dadaxiaoxiao/web-bff/internal/domain"
	"github.com/dadaxiaoxiao/web-bff/internal/job/executor"
	"golang.org/x/sync/semaphore"
	"time"
)

// Scheduler 任务调度
type Scheduler struct {
	execs   map[string]executor.Executor
	client  cronjobv1.CronJobServiceClient
	l       accesslog.Logger
	limiter *semaphore.Weighted
}

func NewScheduler(client cronjobv1.CronJobServiceClient, l accesslog.Logger) *Scheduler {
	return &Scheduler{
		client:  client,
		l:       l,
		limiter: semaphore.NewWeighted(200),
		execs:   make(map[string]executor.Executor),
	}
}

// RegisterExecutor 注册执行器
func (s *Scheduler) RegisterExecutor(exec executor.Executor) {
	s.execs[exec.Name()] = exec
}

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			// 退出调度循环
			return ctx.Err()
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			return err
		}

		// 抢占任务 job
		//  一次抢占任务的数据库查询时间
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		resp, err := s.client.Preempt(dbCtx, &cronjobv1.PreemptRequest{})
		j := resp.Cronjob
		cancel()

		if err != nil {
			s.l.Error("抢占任务失败", accesslog.Error(err))
		}

		// 如何调度执行
		// 找到对应的执行器
		exec, ok := s.execs[j.Executor]
		if !ok {
			// DEBUG 的时候最好中断
			// 线上就继续
			s.l.Error("未找到对应的执行器",
				accesslog.String("executor", j.Executor))
		}

		// 接下来就是执行
		// 单独开启goroutine 来执行 ，不要阻塞主调度循环
		go func() {
			// 最后要释放抢占到任务job
			defer func() {
				s.limiter.Release(1)
				// 释放
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				_, err1 := s.client.Release(ctx, &cronjobv1.ReleaseRequest{Cronjob: j})

				if err1 != nil {
					s.l.Error("释放mysql任务失败",
						accesslog.Error(err1),
						accesslog.Int64("JobId", j.Id))
				}
			}()

			err1 := exec.Exec(ctx, domain.Job{
				Id:       j.Id,
				Cfg:      j.Cfg,
				Executor: j.Executor,
				Name:     j.Name,
				Cron:     j.Expression,
			})

			if err1 != nil {
				s.l.Error("任务执行失败",
					accesslog.Error(err1),
					accesslog.String("executor", j.Executor),
					accesslog.String("name", j.Name))
			}

			// 目前抢占的逻辑
			// 1.没有人调度 tatus= 0  AND next_time <= now
			// 2.曾经有人调度，但是续约失败了，就是这个节点崩溃了

			// status 是通过 CancelFunc 变更 jobStatusWaiting
			// 这里要考虑下一次调度 next_time
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			_, err2 := s.client.ResetNextTime(ctx, &cronjobv1.ResetNextTimeRequest{
				Cronjob: j,
			})
			if err2 != nil {
				s.l.Error("设置下一次执行时间失败", accesslog.Error(err2))
			}
		}()
	}
}
