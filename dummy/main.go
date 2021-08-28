package main

import (
	"time"

	libob "github.com/botuniverse/go-libonebot"
	log "github.com/sirupsen/logrus"
)

type OneBotDummy struct {
	*libob.OneBot
}

func main() {
	log.SetLevel(log.DebugLevel)

	ob := &OneBotDummy{OneBot: libob.NewOneBot("dummy")}

	ob.ActionMux.HandleFunc(libob.ActionGetVersion, func(w libob.ResponseWriter, r *libob.Request) {
		w.WriteData(map[string]string{
			"version":         "1.0.0",
			"onebot_standard": "v12",
		})
	})

	ob.ActionMux.HandleFunc(libob.ActionSendMessage, func(w libob.ResponseWriter, r *libob.Request) {
		userID, err := r.Params.GetString("user_id")
		if err != nil {
			w.WriteFailed(libob.RetCodeParamError, err)
		}
		msg, err := r.Params.GetMessage("message")
		if err != nil {
			w.WriteFailed(libob.RetCodeParamError, err)
		}
		log.Debugf("Send message: %#v, to %v", msg, userID)
	})

	ob.ActionMux.HandleFuncExtended("do_something", func(w libob.ResponseWriter, r *libob.Request) {
	})

	go func() {
		for {
			ob.PushEvent(
				&libob.MessageEvent{
					Event: libob.Event{
						Platform:   "qq",
						SelfID:     "123",
						Type:       libob.EventTypeMessage,
						DetailType: "private",
					},
					UserID:  "234",
					Message: libob.Message{libob.TextSegment("hello")},
				},
			)
			time.Sleep(time.Duration(3) * time.Second)
		}
	}()

	ob.Run()
}
