package onebot

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (ob *OneBot) HandleFunc(action coreAction, handler func(ResponseWriter, *Request)) {
	ob.Handle(action, HandlerFunc(handler))
}

func (ob *OneBot) Handle(action coreAction, handler Handler) {
	ob.handlers[action.string] = handler
}

func (ob *OneBot) HandleFuncExtended(action string, handler func(ResponseWriter, *Request)) {
	ob.HandleExtended(action, HandlerFunc(handler))
}

func (ob *OneBot) HandleExtended(action string, handler HandlerFunc) {
	// if the prefix is empty, then the action name starts with "_"
	ob.extendedHandlers[action] = handler
}

func (ob *OneBot) handleAction(actionBody string) (resp Response) {
	// return "ok" if otherwise explicitly set to "failed"
	w := ResponseWriter{resp: &resp}
	w.WriteOK()

	// try parse the request from the JSON string
	r, err := parseActionRequest(ob.Platform, actionBody)
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
		handlers = &ob.extendedHandlers
	} else {
		handlers = &ob.handlers
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
