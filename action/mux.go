package action

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type Mux struct {
	prefix           string // prefix for extended actions
	handlers         map[string]Handler
	extendedHandlers map[string]Handler
}

func NewMux(prefix string) *Mux {
	return &Mux{
		prefix:           prefix,
		handlers:         map[string]Handler{},
		extendedHandlers: map[string]Handler{},
	}
}

type Handler interface {
	HandleRequest()
}

type HandlerFunc func()

func (handler HandlerFunc) HandleRequest() {
	handler()
}

func (mux *Mux) HandleFunc(action coreAction, handler func()) {
	mux.Handle(action, HandlerFunc(handler))
}

func (mux *Mux) Handle(action coreAction, handler Handler) {
	mux.handlers[action.string] = handler
}

func (mux *Mux) HandleFuncExtended(action string, handler func()) {
	mux.HandleExtended(action, HandlerFunc(handler))
}

func (mux *Mux) HandleExtended(action string, handler HandlerFunc) {
	// if the prefix is empty, then the action name starts with "_"
	mux.extendedHandlers[action] = handler
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

func (mux *Mux) parseRequest(body string) (Request, error) {
	if !gjson.Valid(body) {
		return Request{}, errors.New("Action 请求体不是合法的 JSON")
	}

	bodyJSON := gjson.Parse(body)
	err := validateActionJSON(bodyJSON)
	if err != nil {
		return Request{}, err
	}

	var action Action
	fullname := bodyJSON.Get("action").String()
	if strings.HasPrefix(fullname, mux.prefix+"_") {
		// extended action
		action = Action{
			Prefix: mux.prefix,
			Name:   strings.TrimPrefix(fullname, mux.prefix+"_"),
		}
	} else {
		// core action
		action = Action{
			Prefix: "",
			Name:   fullname,
		}
	}

	r := Request{
		Action: action,
		Params: bodyJSON.Get("params"),
		echo:   bodyJSON.Get("echo"),
	}
	return r, nil
}

func (mux *Mux) HandleRequest(actionBody string) Response {
	log.Debugf("handlers: %#v", mux.handlers)
	log.Debugf("extendedHandlers: %#v", mux.extendedHandlers)

	r, err := mux.parseRequest(actionBody)
	if err != nil {
		errMsg := fmt.Sprintf("Action 请求解析失败: %v", err)
		log.Warnf(errMsg)
		return FailedResponse(RetCodeInvalidRequest, errMsg)
	}

	log.Debugf("Action request: %#v", r)
	// TODO: now it simply return the request
	return Response{
		Status:  StatusOK,
		RetCode: RetCodeOK,
		Data:    r.Params.Value(),
		Message: "",
	}
}
