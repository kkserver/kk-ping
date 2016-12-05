package ping

import (
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
)

type Ping struct {
	Name    string                 `json:"name,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
	Address string                 `json:"address,omitempty"`
	Atime   int64                  `json:"atime,omitempty"`
}

type PingApp struct {
	app.App
	Remote *remote.Service
	Ping   *PingQueryService
}
