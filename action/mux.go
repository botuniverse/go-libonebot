package action

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type ActionMux struct {
	prefix   string // prefix for extended actions
	handlers map[string]Handler
}

func NewActionMux(prefix string) *ActionMux {
	return &ActionMux{
		prefix:   prefix,
		handlers: map[string]Handler{},
	}
}

type Handler interface {
	HandleRequest()
}

type HandlerFunc func()

func (handler HandlerFunc) HandleRequest() {
	handler()
}

func (mux *ActionMux) HandleFunc(action coreAction, handler func()) {
	mux.Handle(action, HandlerFunc(handler))
}

func (mux *ActionMux) Handle(action coreAction, handler Handler) {
	mux.handlers[action.string] = handler
}

func (mux *ActionMux) HandleFuncExtended(action string, handler func()) {
	mux.HandleExtended(action, HandlerFunc(handler))
}

func (mux *ActionMux) HandleExtended(action string, handler HandlerFunc) {
	// if the prefix is empty, then the action name starts with "_"
	mux.handlers[mux.prefix+"_"+action] = handler
}

func validateActionJSON(actionJSON gjson.Result) error {
	if !actionJSON.Get("action").Exists() {
		return errors.New("Action 请求体缺少 `action` 字段")
	}
	if actionJSON.Get("action").String() == "" {
		return errors.New("Action 请求体的 `action` 字段为空")
	}
	if !actionJSON.Get("params").Exists() {
		return errors.New("Action 请求体缺少 `params` 字段")
	}
	if !actionJSON.Get("params").IsObject() {
		return errors.New("Action 请求体的 `params` 字段不是一个 JSON 对象")
	}
	return nil
}

func (mux *ActionMux) ParseRequest(actionBody string) (Request, error) {
	if !gjson.Valid(actionBody) {
		return Request{}, errors.New("Action 请求体不是合法的 JSON")
	}

	actionJSON := gjson.Parse(actionBody)
	err := validateActionJSON(actionJSON)
	if err != nil {
		return Request{}, err
	}

	var action Action
	actionFullname := actionJSON.Get("action").String()
	if strings.HasPrefix(actionFullname, mux.prefix+"_") {
		action = Action{
			Prefix: mux.prefix,
			Name:   strings.TrimPrefix(actionFullname, mux.prefix+"_"),
		}
	} else {
		action = Action{
			Prefix: "",
			Name:   actionFullname,
		}
	}

	r := Request{
		Action: action,
		Params: actionJSON.Get("params"),
		echo:   actionJSON.Get("echo"),
	}
	return r, nil
}

// TODO: input and output types
func (mux *ActionMux) HandleRequest(r *Request) Response {
	log.Debugf("handlers: %#v", mux.handlers)
	log.Debugf("Action request: %#v", r)
	// TODO: now it simply return the request
	return Response{
		Status:  StatusOK,
		RetCode: RetCodeOK,
		Data:    r.Params.Value(),
		Message: "",
	}
}
