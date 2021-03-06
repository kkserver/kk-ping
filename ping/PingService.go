package ping

import (
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
	"github.com/kkserver/kk-lib/kk/json"
	"log"
	"strings"
	"time"
)

type PingService struct {
	app.Service

	ReceiveMessage *remote.RemoteReceiveMessageTask
	Query          *PingQueryTask

	Expires int64

	dispatch *kk.Dispatch
	pings    map[int64]*Ping
	id       int64
}

func (S *PingService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *PingService) HandleInitTask(a app.IApp, task *app.InitTask) error {

	S.dispatch = kk.NewDispatch()
	S.pings = map[int64]*Ping{}
	S.id = 0

	var fn func() = nil

	fn = func() {

		var ids = []int64{}
		var now = time.Now().Unix()

		for id, ping := range S.pings {

			if ping.Status == PingStatusOffline || ping.Atime+S.Expires < now {
				ids = append(ids, id)
			}

		}

		for _, id := range ids {
			delete(S.pings, id)
		}

		log.Println("PingService", "Cleanup", ids)

		S.dispatch.AsyncDelay(fn, 30*time.Second)
	}

	S.dispatch.Async(fn)

	return nil
}

func (S *PingService) HandleRemoteReceiveMessageTask(a app.IApp, task *remote.RemoteReceiveMessageTask) error {

	log.Println(task.Message.String())

	if task.Message.Method == "PING" {

		var v = Ping{}

		err := json.Decode(task.Message.Content, &v)

		if err == nil {

			v.Name = task.Message.From
			v.Atime = time.Now().Unix()
			v.Status = PingStatusOnline

			S.dispatch.Async(func() {

				for _, ping := range S.pings {

					if ping.Name == v.Name && ping.Address == v.Address {
						ping.Atime = v.Atime
						ping.Status = v.Status
						if v.Options != nil {
							ping.Options = v.Options
						}
						if v.Counter != nil {
							ping.Counter = v.Counter
						}
						if v.Tasks != nil {
							ping.Tasks = v.Tasks
						}
						return
					}

				}

				if v.Ctime == 0 {
					v.Ctime = v.Atime
				}

				S.pings[S.id] = &v

				S.id = S.id + 1

			})

		} else {
			log.Println("[PingService] " + err.Error())
		}

	} else if task.Message.Method == "PING.DISCONNECTED" {

		var v = Ping{}

		err := json.Decode(task.Message.Content, &v)

		if err == nil {

			v.Name = task.Message.From
			v.Atime = time.Now().Unix()
			v.Status = PingStatusOffline

			S.dispatch.Async(func() {

				for _, ping := range S.pings {

					if ping.Name == v.Name && ping.Address == v.Address {
						ping.Atime = v.Atime
						ping.Status = v.Status
						if v.Options != nil {
							ping.Options = v.Options
						}
						if v.Counter != nil {
							ping.Counter = v.Counter
						}
						if v.Tasks != nil {
							ping.Tasks = v.Tasks
						}
						return
					}

				}

			})

		} else {
			log.Println("[PingService] " + err.Error())
		}

	}

	return nil
}

func (S *PingService) HandlePingQueryTask(a app.IApp, task *PingQueryTask) error {

	S.dispatch.Sync(func() {

		var pings = []*Ping{}

		for _, ping := range S.pings {

			if (task.Name == "" || ping.Name == task.Name) && (task.Address == "" || ping.Address == task.Address) &&
				(task.Prefix == "" || strings.HasPrefix(ping.Name, task.Prefix)) {
				pings = append(pings, ping)
			}

		}

		task.Result.Pings = pings

	})

	return nil
}
