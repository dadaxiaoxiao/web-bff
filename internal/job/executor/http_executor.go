package executor

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dadaxiaoxiao/web-bff/internal/domain"
	"net/http"
)

// HttpExecutor http 执行器
type HttpExecutor struct {
}

func (h HttpExecutor) Name() string {
	return "http"
}

func (h HttpExecutor) Exec(ctx context.Context, j domain.Job) error {
	type Config struct {
		EndPoint string
		Method   string
	}
	var cfg Config
	err := json.Unmarshal([]byte(j.Cfg), &cfg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(cfg.Method, cfg.EndPoint, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		return errors.New("执行失败")
	}
	return nil
}

func NewHttpExecutor() *HttpExecutor {
	return &HttpExecutor{}
}
