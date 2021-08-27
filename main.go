package main

import (
	"time"

	"github.com/botuniverse/go-libonebot/comm"
	"github.com/botuniverse/go-libonebot/event"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)

	eventDispatcher := event.NewEventDispatcher()
	comm.StartHTTPTask("127.0.0.1", 5700)
	comm.StartWSTask("127.0.0.1", 6700, eventDispatcher)
	comm.StartHTTPWebhookTask("http://127.0.0.1:8080", eventDispatcher)

	time.Sleep(time.Duration(3) * time.Second)
	eventDispatcher.Dispatch(
		&event.MessageEvent{
			Event: event.Event{
				Platform:   "qq",
				SelfID:     "123",
				Type:       event.TypeMessage,
				DetailType: "private",
			},
			UserID:  "234",
			Message: "hello",
		},
	)

	select {}
}
