package ping

import (
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
	"github.com/kkserver/kk-lib/kk/json"
	"log"
	"time"
)

type PingQueryService struct {
	app.Service
	Init           *app.InitTask
	ReceiveMessage *remote.RemoteReceiveMessageTask
	Query          *PingQueryTask
	Cleanup        *PingCleanupTask

	dispatch *kk.Dispatch
	pings    map[int64]Ping
	id       int64
}

func (S *PingQueryService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *PingQueryService) HandleInitTask(a app.IApp, task *app.InitTask) error {

	S.dispatch = kk.NewDispatch()
	S.pings = map[int64]Ping{}
	S.id = 0

	return nil
}

func (S *PingQueryService) HandleRemoteReceiveMessageTask(a app.IApp, task *remote.RemoteReceiveMessageTask) error {

	if task.Message.Method == "PING" {

		var v = Ping{}

		err := json.Decode(task.Message.Content, &v)

		if err == nil {

			v.Name = task.Message.From
			v.Atime = time.Now().Unix()

			S.dispatch.Async(func() {

				for _, ping := range S.pings {

					if ping.Name == v.Name && ping.Address == v.Address {
						ping.Options = v.Options
						ping.Atime = v.Atime
						return
					}

				}

				S.pings[S.id] = v

				S.id = S.id + 1

			})

		} else {
			log.Println("[PingService] " + err.Error())
		}

	}

	return nil
}

func (S *PingQueryService) HandlePingQueryTask(a app.IApp, task *PingQueryTask) error {

	S.dispatch.Sync(func() {

		var pings = []Ping{}

		for _, ping := range S.pings {

			if (task.Name == "" || ping.Name == task.Name) && (task.Address == "" || ping.Address == task.Address) {
				pings = append(pings, ping)
			}

		}

		task.Result.Pings = pings

	})

	return nil
}

func (S *PingQueryService) HandlePingCleanupTask(a app.IApp, task *PingCleanupTask) error {

	S.dispatch.Sync(func() {

		var ids = []int64{}
		var now = time.Now().Unix()

		for id, ping := range S.pings {

			if ping.Atime+task.Expires < now {
				ids = append(ids, id)
			}

		}

		for _, id := range ids {
			delete(S.pings, id)
		}

		log.Println("[PingQueryService][HandlePingCleanupTask] ", ids)

	})

	return nil
}
