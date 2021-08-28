package libonebot

import (
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
)

type EventBroadcaster struct {
	listenChans     []chan []byte
	listenChansLock sync.RWMutex
}

func NewEventBroadcaster() *EventBroadcaster {
	return &EventBroadcaster{
		listenChans: make([]chan []byte, 0),
	}
}

func (broadcaster *EventBroadcaster) OpenListenChan() <-chan []byte {
	broadcaster.listenChansLock.Lock()
	defer broadcaster.listenChansLock.Unlock()

	ch := make(chan []byte) // TODO: channel size
	broadcaster.listenChans = append(broadcaster.listenChans, ch)
	return ch
}

func (broadcaster *EventBroadcaster) CloseListenChan(listenCh <-chan []byte) {
	broadcaster.listenChansLock.Lock()
	defer broadcaster.listenChansLock.Unlock()

	for i, ch := range broadcaster.listenChans {
		if ch == listenCh {
			close(ch)
			broadcaster.listenChans = append(broadcaster.listenChans[:i], broadcaster.listenChans[i+1:]...)
			return
		}
	}
}

func (broadcaster *EventBroadcaster) Broadcast(event AnyEvent) bool {
	log.Debugf("Event: %#v", event)

	if !event.TryFixUp() {
		log.Warnf("事件字段值无效")
		return false
	}

	jsonBytes, err := json.Marshal(event)
	if err != nil {
		log.Warnf("事件序列化失败, 错误: %v", err)
		return false
	}

	broadcaster.listenChansLock.RLock() // use read lock to allow emitting events concurrently
	defer broadcaster.listenChansLock.RUnlock()
	for _, ch := range broadcaster.listenChans {
		ch <- jsonBytes
	}
	return true
}
