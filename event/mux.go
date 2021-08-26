package event

type EventMux struct {
	outChannels []chan *Event
}

func NewEventMux() *EventMux {
	return &EventMux{
		outChannels: make([]chan *Event, 0),
	}
}

func (mux *EventMux) OpenOutChan() <-chan *Event {
	outCh := make(chan *Event) // TODO: channel size
	mux.outChannels = append(mux.outChannels, outCh)
	return outCh
}

func (mux *EventMux) CloseOutChan(outCh <-chan *Event) {
	for i, ch := range mux.outChannels {
		if ch == outCh {
			close(ch)
			mux.outChannels = append(mux.outChannels[:i], mux.outChannels[i+1:]...)
			return
		}
	}
}

func (mux *EventMux) Emit(event *Event) {
	for _, outCh := range mux.outChannels {
		outCh <- event
	}
}
