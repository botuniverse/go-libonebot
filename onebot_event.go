package onebot

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func (ob *OneBot) PushEvent(event AnyEvent) bool {
	log.Debugf("Event: %#v", event)

	if !event.tryFixUp() {
		log.Warnf("事件字段值无效")
		return false
	}

	jsonBytes, err := json.Marshal(event)
	if err != nil {
		log.Warnf("事件序列化失败, 错误: %v", err)
		return false
	}

	ob.eventListenChansLock.RLock() // use read lock to allow emitting events concurrently
	defer ob.eventListenChansLock.RUnlock()
	for _, ch := range ob.eventListenChans {
		ch <- jsonBytes
	}
	return true
}

func (ob *OneBot) openEventListenChan() <-chan []byte {
	ob.eventListenChansLock.Lock()
	defer ob.eventListenChansLock.Unlock()

	ch := make(chan []byte) // TODO: channel size
	ob.eventListenChans = append(ob.eventListenChans, ch)
	return ch
}

func (ob *OneBot) closeEventListenChan(ch <-chan []byte) {
	ob.eventListenChansLock.Lock()
	defer ob.eventListenChansLock.Unlock()

	for i, c := range ob.eventListenChans {
		if c == ch {
			ob.eventListenChans = append(ob.eventListenChans[:i], ob.eventListenChans[i+1:]...)
			return
		}
	}
}
