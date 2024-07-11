package job

import "github.com/dadaxiaoxiao/web-bff/internal/job/executor"

// Scheduler 任务调度
type Scheduler struct {
	execs map[string]executor.Executor
}
