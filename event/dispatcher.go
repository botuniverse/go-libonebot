package event

import (
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
)

type EventDispatcher struct {
	outChans     []chan []byte
	outChansLock sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		outChans: make([]chan []byte, 0),
	}
}

func (dispatcher *EventDispatcher) OpenOutChan() <-chan []byte {
	dispatcher.outChansLock.Lock()
	defer dispatcher.outChansLock.Unlock()

	outCh := make(chan []byte) // TODO: channel size
	dispatcher.outChans = append(dispatcher.outChans, outCh)
	return outCh
}

func (dispatcher *EventDispatcher) CloseOutChan(outCh <-chan []byte) {
	dispatcher.outChansLock.Lock()
	defer dispatcher.outChansLock.Unlock()

	for i, ch := range dispatcher.outChans {
		if ch == outCh {
			close(ch)
			dispatcher.outChans = append(dispatcher.outChans[:i], dispatcher.outChans[i+1:]...)
			return
		}
	}
}

func (dispatcher *EventDispatcher) Dispatch(event anyEvent) {
	log.Debugf("Event: %#v", event)
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		log.Warnf("事件序列化失败, 错误: %v", err)
		return
	}

	dispatcher.outChansLock.RLock() // use read lock to allow emitting events concurrently
	defer dispatcher.outChansLock.RUnlock()
	for _, ch := range dispatcher.outChans {
		ch <- jsonBytes
	}
}
