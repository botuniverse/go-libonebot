package main

import (
	"github.com/botuniverse/go-libonebot/action"
	"github.com/botuniverse/go-libonebot/comm"
	"github.com/botuniverse/go-libonebot/event"
	log "github.com/sirupsen/logrus"
)

type OneBot struct {
	Platform        string
	ActionMux       *action.Mux
	eventDispatcher *event.Dispatcher
}

func NewOneBot(platform string) *OneBot {
	if platform == "" {
		log.Warnf("没有设置 OneBot 实现平台名称, 可能导致程序行为与预期不符")
	}
	return &OneBot{
		Platform:        platform,
		ActionMux:       action.NewMux(platform),
		eventDispatcher: event.NewDispatcher(),
	}
}

func (ob *OneBot) startCommunicationMethods() {
	comm.StartHTTPTask("127.0.0.1", 5700, ob.ActionMux)
	comm.StartWSTask("127.0.0.1", 6700, ob.ActionMux, ob.eventDispatcher)
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
