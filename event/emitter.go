package event

import (
	"encoding/json"
	"sync"

	log "github.com/sirupsen/logrus"
)

type EventEmitter struct {
	outChans     []chan []byte
	outChansLock sync.RWMutex
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		outChans: make([]chan []byte, 0),
	}
}

func (emitter *EventEmitter) OpenOutChan() <-chan []byte {
	emitter.outChansLock.Lock()
	defer emitter.outChansLock.Unlock()

	outCh := make(chan []byte) // TODO: channel size
	emitter.outChans = append(emitter.outChans, outCh)
	return outCh
}

func (emitter *EventEmitter) CloseOutChan(outCh <-chan []byte) {
	emitter.outChansLock.Lock()
	defer emitter.outChansLock.Unlock()

	for i, ch := range emitter.outChans {
		if ch == outCh {
			close(ch)
			emitter.outChans = append(emitter.outChans[:i], emitter.outChans[i+1:]...)
			return
		}
	}
}

func (emitter *EventEmitter) Emit(event anyEvent) {
	log.Debugf("Event: %#v", event)
	jsonBytes, err := json.Marshal(event)
	if err != nil {
		log.Warnf("事件序列化失败, 错误: %v", err)
		return
	}

	emitter.outChansLock.RLock() // use read lock to allow emitting events concurrently
	defer emitter.outChansLock.RUnlock()
	for _, ch := range emitter.outChans {
		ch <- jsonBytes
	}
}
