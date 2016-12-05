package ping

import (
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
	"github.com/kkserver/kk-lib/kk/json"
	"log"
	"time"
)

type PingService struct {
	app.Service
	Options        map[string]interface{}
	Expires        int64
	ToName         string
	ReceiveMessage *remote.RemoteReceiveMessageTask

	address string
}

func (S *PingService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *PingService) Ping(a app.IApp) error {

	task := remote.RemoteSendMessageTask{}

	task.Message.Method = "PING"
	task.Message.To = S.ToName
	task.Message.Type = "text/json"

	var v = Ping{}

	v.Options = S.Options
	v.Address = S.address

	task.Message.Content, _ = json.Encode(&v)

	kk.GetDispatchMain().AsyncDelay(func() {
		S.Ping(a)
	}, time.Duration(S.Expires)*time.Second)

	log.Println("PING")
	log.Println(v)

	return app.Handle(a, &task)
}

func (S *PingService) HandleRemoteReceiveMessageTask(a app.IApp, task *remote.RemoteReceiveMessageTask) error {

	if task.Message.Method == "CONNECTED" {

		S.address = string(task.Message.Content)

		kk.GetDispatchMain().Async(func() {
			S.Ping(a)
		})

	}

	return nil
}
