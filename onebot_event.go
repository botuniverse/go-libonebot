package libonebot

import "encoding/json"

// Push 向与 OneBot 实例连接的接受端推送一个事件.
func (ob *OneBot) Push(event AnyEvent) bool {
	return ob.PushWithSelf(event, ob.Self)
}

// PushWithSelf 向与 OneBot 实例连接的接受端推送一个事件, 并指定收到事件的机器人自身标识.
func (ob *OneBot) PushWithSelf(event AnyEvent, self *Self) bool {
	if event == nil {
		ob.Logger.Errorf("事件为空")
		return false
	}
	if err := event.tryFixUp(self); err != nil {
		ob.Logger.Errorf("事件无效, 错误: %v", err)
		return false
	}
	ob.Logger.Debugf("事件: %#v", event)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		ob.Logger.Errorf("事件序列化失败, 错误: %v", err)
		return false
	}

	ob.Logger.Infof("事件 `%v` 开始推送", event.Name())
	ob.eventListenChansLock.RLock() // use read lock to allow emitting events concurrently
	defer ob.eventListenChansLock.RUnlock()
	for _, ch := range ob.eventListenChans {
		ch <- marshaledEvent{event.Name(), eventBytes, event}
	}
	return true
}

type marshaledEvent struct {
	name  string
	bytes []byte
	raw   AnyEvent
}

func (ob *OneBot) openEventListenChan() <-chan marshaledEvent {
	ch := make(chan marshaledEvent, 1) //  设置缓冲区为1，因为需要先放入一个连接事件
	connectMetaEvent := MakeConnectMetaEvent(ob.Impl, Version, OneBotVersion)
	ob.Logger.Debugf("事件: %#v", connectMetaEvent)
	ob.Logger.Infof("事件 `%v` 开始推送", connectMetaEvent.Name())
	eventBytes, _ := json.Marshal(connectMetaEvent)
	ch <- marshaledEvent{
		name:  connectMetaEvent.Name(),
		bytes: eventBytes,
		raw:   &connectMetaEvent,
	}
	ob.eventListenChansLock.Lock()
	ob.eventListenChans = append(ob.eventListenChans, ch)
	ob.eventListenChansLock.Unlock()
	return ch
}

func (ob *OneBot) closeEventListenChan(ch <-chan marshaledEvent) {
	ob.eventListenChansLock.Lock()
	defer ob.eventListenChansLock.Unlock()

	for i, c := range ob.eventListenChans {
		if c == ch {
			close(c)
			ob.eventListenChans = append(ob.eventListenChans[:i], ob.eventListenChans[i+1:]...)
			return
		}
	}
}
