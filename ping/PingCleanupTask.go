package ping

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type PingCleanupTaskResult struct {
	app.Result
}

type PingCleanupTask struct {
	app.Task
	Expires int64 `json:"expires,omitempty"`
	Result  PingCleanupTaskResult
}

func (task *PingCleanupTask) GetResult() interface{} {
	return &task.Result
}

func (task *PingCleanupTask) GetInhertType() string {
	return "ping"
}

func (task *PingCleanupTask) GetClientName() string {
	return "Ping.Cleanup"
}
