package main

import (
	"github.com/dadaxiaoxiao/go-pkg/customserver"
	"github.com/dadaxiaoxiao/web-bff/internal/job"
)

type WebApp struct {
	customserver.App
	Scheduler *job.Scheduler
}
