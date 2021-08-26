package event

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type EventEmitter struct {
	outChans []chan []byte
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		outChans: make([]chan []byte, 0),
	}
}

func (emitter *EventEmitter) OpenOutChan() <-chan []byte {
	outCh := make(chan []byte) // TODO: channel size
	emitter.outChans = append(emitter.outChans, outCh)
	return outCh
}

func (emitter *EventEmitter) CloseOutChan(outCh <-chan []byte) {
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
	for _, ch := range emitter.outChans {
		ch <- jsonBytes
	}
}
