package libonebot

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type ActionMux struct {
	prefix           string // prefix for extended actions
	handlers         map[string]Handler
	extendedHandlers map[string]Handler
}

func NewActionMux(prefix string) *ActionMux {
	return &ActionMux{
		prefix:           prefix,
		handlers:         map[string]Handler{},
		extendedHandlers: map[string]Handler{},
	}
}

func (handler HandlerFunc) HandleAction(w ResponseWriter, r *Request) {
	handler(w, r)
}

func (mux *ActionMux) HandleFunc(action coreAction, handler func(ResponseWriter, *Request)) {
	mux.Handle(action, HandlerFunc(handler))
}

func (mux *ActionMux) Handle(action coreAction, handler Handler) {
	mux.handlers[action.string] = handler
}

func (mux *ActionMux) HandleFuncExtended(action string, handler func(ResponseWriter, *Request)) {
	mux.HandleExtended(action, HandlerFunc(handler))
}

func (mux *ActionMux) HandleExtended(action string, handler HandlerFunc) {
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

func (mux *ActionMux) parseRequest(body string) (Request, error) {
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
	prefix := mux.prefix + "_"
	if strings.HasPrefix(fullname, prefix) {
		// extended action
		action = Action{
			Prefix:     mux.prefix,
			Name:       strings.TrimPrefix(fullname, prefix),
			IsExtended: true,
		}
	} else {
		// core action
		action = Action{
			Prefix:     "",
			Name:       fullname,
			IsExtended: false,
		}
	}

	r := Request{
		Action: action,
		Params: Params{JSON: bodyJSON.Get("params")},
		Echo:   bodyJSON.Get("echo").Value(),
	}
	return r, nil
}

func (mux *ActionMux) HandleAction(actionBody string) (resp Response) {
	// return "ok" if otherwise explicitly set to "failed"
	w := ResponseWriter{resp: &resp}
	w.WriteOK()

	// try parse the request from the JSON string
	r, err := mux.parseRequest(actionBody)
	if err != nil {
		err := fmt.Errorf("Action 请求解析失败, 错误: %v", err)
		log.Warn(err)
		w.WriteFailed(RetCodeInvalidRequest, err)
		return
	}
	log.Debugf("Action request: %#v", r)

	// once we got the `echo` field, set the `echo` field in the response
	resp.Echo = r.Echo

	var handlers *map[string]Handler
	if r.Action.IsExtended {
		handlers = &mux.extendedHandlers
	} else {
		handlers = &mux.handlers
	}

	handler := (*handlers)[r.Action.Name]
	if handler == nil {
		err := fmt.Errorf("Action `%v` 不存在", r.Action)
		log.Warn(err)
		w.WriteFailed(RetCodeActionNotFound, err)
		return
	}

	log.Infof("Action `%v` 开始处理", r.Action)
	handler.HandleAction(w, &r)
	return
}
