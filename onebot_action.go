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

func (ob *OneBot) handleAction(r *Request) (resp Response) {
	ob.Logger.Debugf("动作请求: %+v", r)
	resp.Echo = r.Echo
	w := ResponseWriter{resp: &resp}

	if ob.actionHandler == nil {
		err := fmt.Errorf("动作请求处理器未设置")
		ob.Logger.Warn(err)
		w.WriteFailed(RetCodeActionNotFound, err)
		return
	}

	ob.Logger.Debugf("动作请求 `%v` 开始处理", r.Action)
	ob.actionHandler.HandleAction(w, r)
	if resp.Status.string == "" {
		err := fmt.Errorf("动作请求处理器没有正确设置响应状态")
		ob.Logger.Warn(err)
		w.WriteFailed(RetCodeBadActionHandler, err)
		return
	}
	if resp.Status != statusOK {
		ob.Logger.Warnf("动作请求 `%v` 处理失败, 错误: %v", r.Action, resp.Message)
	} else {
		ob.Logger.Infof("动作请求 `%v` 处理成功", r.Action)
	}
	return
}

func (ob *OneBot) parseAction(actionBytes []byte, isBinary bool) (Request, error) {
	if isBinary {
		return parseBinaryActionRequest(ob.Platform, actionBytes)
	}
	return parseTextActionRequest(ob.Platform, actionBytes)
}

func (ob *OneBot) parseAndHandleAction(actionBytes []byte, isBinary bool) Response {
	request, err := ob.parseAction(actionBytes, isBinary)
	if err != nil {
		err := fmt.Errorf("动作请求解析失败, 错误: %v", err)
		ob.Logger.Warn(err)
		return failedResponse(RetCodeInvalidRequest, err)
	}
	return ob.handleAction(&request)
}
