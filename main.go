package main

import (
	"time"

	"github.com/botuniverse/go-libonebot/action"
	"github.com/botuniverse/go-libonebot/event"
	"github.com/botuniverse/go-libonebot/message"
	log "github.com/sirupsen/logrus"
)

type OneBotDummy struct {
	*OneBot
}

func main() {
	log.SetLevel(log.DebugLevel)

	ob := &OneBotDummy{OneBot: NewOneBot("dummy")}

	ob.ActionMux.HandleFunc(action.ActionGetVersion, func(w action.ResponseWriter, r *action.Request) {
		w.WriteData(map[string]string{
			"version":         "1.0.0",
			"onebot_standard": "v12",
		})
	})

	ob.ActionMux.HandleFunc(action.ActionSendMessage, func(w action.ResponseWriter, r *action.Request) {
		userID, err := r.Params.GetString("user_id")
		if err != nil {
			w.WriteFailed(action.RetCodeParamError, err.Error()) // TODO
		}
		msg, err := r.Params.GetMessage("message")
		if err != nil {
			w.WriteFailed(action.RetCodeParamError, err.Error())
		}
		log.Debugf("Send message: %#v, to %v", msg, userID)
	})

	ob.ActionMux.HandleFuncExtended("do_something", func(w action.ResponseWriter, r *action.Request) {
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
					Message: message.Message{message.TextSegment("hello")},
				},
			)
			time.Sleep(time.Duration(3) * time.Second)
		}
	}()

	ob.Run()
}
