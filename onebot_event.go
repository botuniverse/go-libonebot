package libonebot

import (
	"encoding/json"
)

// Push 向与 OneBot 实例连接的接受端推送一个事件.
func (ob *OneBot) Push(event AnyEvent) bool {
	if event == nil {
		ob.Logger.Warnf("事件为空")
		return false
	}
	if !event.tryFixUp(ob.Platform, ob.SelfID) {
		ob.Logger.Warnf("事件字段值无效")
		return false
	}
	ob.Logger.Debugf("事件: %#v", event)

	eventJSONBytes, err := json.Marshal(event)
	if err != nil {
		ob.Logger.Warnf("事件序列化失败, 错误: %v", err)
		return false
	}

	ob.Logger.Infof("事件 `%v` 开始推送", event.Name())
	ob.eventListenChansLock.RLock() // use read lock to allow emitting events concurrently
	defer ob.eventListenChansLock.RUnlock()
	for _, ch := range ob.eventListenChans {
		ch <- marshaledEvent{event.Name(), eventJSONBytes, event}
	}
	return true
}

type marshaledEvent struct {
	name  string
	bytes []byte
	raw   AnyEvent
}

func (ob *OneBot) openEventListenChan() <-chan marshaledEvent {
	ch := make(chan marshaledEvent) // TODO: channel size
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
