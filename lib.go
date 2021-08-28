package libonebot

import (
	log "github.com/sirupsen/logrus"
)

type OneBot struct {
	Platform         string
	ActionMux        *ActionMux
	eventBroadcaster *eventBroadcaster
}

func NewOneBot(platform string) *OneBot {
	if platform == "" {
		log.Warnf("没有设置 OneBot 实现平台名称, 可能导致程序行为与预期不符")
	}
	return &OneBot{
		Platform:         platform,
		ActionMux:        NewActionMux(platform),
		eventBroadcaster: newEventBroadcaster(),
	}
}

func (ob *OneBot) startCommunicationMethods() {
	commStartHTTP("127.0.0.1", 5700, ob.ActionMux)
	commStartWS("127.0.0.1", 6700, ob.ActionMux, ob.eventBroadcaster)
	commStartHTTPWebhook("http://127.0.0.1:8080", ob.eventBroadcaster)
}

func (ob *OneBot) Run() {
	ob.startCommunicationMethods()
	log.Infof("OneBot 运行中...")
	select {}
}

func (ob *OneBot) PushEvent(event AnyEvent) {
	ob.eventBroadcaster.Broadcast(event)
}
