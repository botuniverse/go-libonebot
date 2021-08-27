package main

import (
	"github.com/botuniverse/go-libonebot/comm"
	"github.com/botuniverse/go-libonebot/event"
	log "github.com/sirupsen/logrus"
)

type OneBot struct {
	eventDispatcher *event.EventDispatcher
}

func NewOneBot() *OneBot {
	return &OneBot{
		eventDispatcher: event.NewEventDispatcher(),
	}
}

func (ob *OneBot) startCommunicationMethods() {
	comm.StartHTTPTask("127.0.0.1", 5700)
	comm.StartWSTask("127.0.0.1", 6700, ob.eventDispatcher)
	comm.StartHTTPWebhookTask("http://127.0.0.1:8080", ob.eventDispatcher)
}

func (ob *OneBot) Run() {
	ob.startCommunicationMethods()
	log.Infof("OneBot 运行中...")
	select {}
}

func (ob *OneBot) PushEvent(event event.AnyEvent) {
	ob.eventDispatcher.Dispatch(event)
}
