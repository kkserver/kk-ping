package ping

import (
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
)

const PingStatusOnline = 200
const PingStatusOffline = 300

type Ping struct {
	Name    string                     `json:"name,omitempty"`
	Options map[string]interface{}     `json:"options,omitempty"`
	Address string                     `json:"address,omitempty"`
	Counter *remote.Counter            `json:"counter,omitempty"`
	Tasks   map[string]*remote.Counter `json:"tasks,omitempty"`
	Status  int                        `json:"status,omitempty"`
	Atime   int64                      `json:"atime,omitempty"`
}

type PingApp struct {
	app.App
	Remote *remote.Service
	Ping   *PingService
}
