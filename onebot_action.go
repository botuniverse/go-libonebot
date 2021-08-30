package libonebot

import (
	"fmt"
)

func (ob *OneBot) HandleFunc(handler func(ResponseWriter, *Request)) {
	ob.Handle(HandlerFunc(handler))
}

func (ob *OneBot) Handle(handler Handler) {
	ob.actionHandler = handler
}

func (ob *OneBot) handleAction(actionBody string) (resp Response) {
	w := ResponseWriter{resp: &resp}

	// try parse the request from the JSON string
	r, err := parseActionRequest(ob.Platform, actionBody)
	if err != nil {
		err := fmt.Errorf("动作请求解析失败, 错误: %v", err)
		ob.Logger.Warn(err)
		w.WriteFailed(RetCodeInvalidRequest, err)
		return
	}
	ob.Logger.Debugf("动作请求: %#v", r)

	// once we got the `echo` field, set the `echo` field in the response
	resp.Echo = r.Echo

	if ob.actionHandler == nil {
		err := fmt.Errorf("动作请求处理器未设置")
		ob.Logger.Warn(err)
		w.WriteFailed(RetCodeActionNotFound, err)
		return
	}

	ob.Logger.Debugf("动作请求 `%v` 开始处理", r.Action)
	ob.actionHandler.HandleAction(w, &r)
	if resp.Status != statusOK {
		ob.Logger.Warnf("动作请求 `%v` 处理失败, 错误: %v", r.Action, resp.Message)
	} else {
		ob.Logger.Infof("动作请求 `%v` 处理成功", r.Action)
	}
	return
}
