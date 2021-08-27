package main

import (
	"time"

	"github.com/botuniverse/go-libonebot/action"
	"github.com/botuniverse/go-libonebot/event"
	log "github.com/sirupsen/logrus"
)

type OneBotDummy struct {
	*OneBot
}

func main() {
	log.SetLevel(log.DebugLevel)

	ob := &OneBotDummy{OneBot: NewOneBot("dummy")}

	ob.ActionMux.HandleFunc(action.ActionGetVersion, func() {
		log.Debugf("ActionGetVersion")
	})

	go func() {
		for {
			ob.PushEvent(
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
			time.Sleep(time.Duration(3) * time.Second)
		}
	}()

	ob.Run()
}
