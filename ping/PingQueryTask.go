package ping

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type PingQueryTaskResult struct {
	app.Result
	Pings []*Ping `json:"pings,omitempty"`
}

type PingQueryTask struct {
	app.Task
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
	Prefix  string `json:"prefix,omitempty"`
	Result  PingQueryTaskResult
}

func (task *PingQueryTask) GetResult() interface{} {
	return &task.Result
}

func (task *PingQueryTask) GetInhertType() string {
	return "ping"
}

func (task *PingQueryTask) GetClientName() string {
	return "Ping.Query"
}
