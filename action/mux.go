package action

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type ActionMux struct {
	prefix   string // prefix for extended actions
	handlers map[string]Handler
}

func NewActionMux(prefix string) *ActionMux {
	return &ActionMux{prefix: prefix}
}

type Handler interface {
	HandleAction()
}

type HandlerFunc func()

func (handler HandlerFunc) HandleAction() {
	handler()
}

func (mux *ActionMux) HandleFunc(action action, handler func()) {
	mux.Handle(action, HandlerFunc(handler))
}

func (mux *ActionMux) Handle(action action, handler Handler) {
	mux.handlers[action.string] = handler
}

func (mux *ActionMux) HandleFuncExtended(action string, handler func()) {
	mux.HandleExtended(action, HandlerFunc(handler))
}

func (mux *ActionMux) HandleExtended(action string, handler HandlerFunc) {
	// if the prefix is empty, then the action name starts with "_"
	mux.handlers[mux.prefix+"_"+action] = handler
}

// TODO: input and output types
func HandleAction(request gjson.Result) gjson.Result {
	log.Debugf("Action request: %v", request)
	// TODO: now it simply return the request
	return request
}
