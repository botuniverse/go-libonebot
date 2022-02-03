package libonebot

import (
	"fmt"
)

// HandleFunc 将一个函数注册为动作处理器.
//
// 一个 OneBot 实例只能注册一个动作处理器, 多次调用将覆盖.
func (ob *OneBot) HandleFunc(handler func(ResponseWriter, *Request)) {
	ob.Handle(HandlerFunc(handler))
}

// Handle 将一个 Handler 对象注册为动作处理器.
//
// 一个 OneBot 实例只能注册一个动作处理器, 多次调用将覆盖.
// 可以传入 ActionMux 对象来根据动作名称分发请求到不同的 Handler 对象.
func (ob *OneBot) Handle(handler Handler) {
	ob.actionHandler = handler
}

// CallAction 调用指定动作.
//
// 参数:
//   action: 要调用的动作名称
//   params: 动作参数, 若传入 nil 则实际动作参数为空 map
func (ob *OneBot) CallAction(action string, params map[string]interface{}) Response {
	if params == nil {
		params = make(map[string]interface{})
	}
	req := &Request{
		Action: action,
		Params: EasierMapFromMap(params),
	}
	return ob.handleRequest(req)
}

func (ob *OneBot) handleRequest(r *Request) (resp Response) {
	ob.Logger.Debugf("动作请求: %+v", r)
	resp.Echo = r.Echo
	w := ResponseWriter{resp: &resp}

	if ob.actionHandler == nil {
		err := fmt.Errorf("动作处理器未设置")
		ob.Logger.Warn(err)
		w.WriteFailed(RetCodeUnsupportedAction, err)
		return
	}

	ob.Logger.Debugf("动作请求 `%v` 开始处理", r.Action)
	ob.actionHandler.HandleAction(w, r)
	if resp.Status == statusOK {
		ob.Logger.Infof("动作请求 `%v` 处理成功", r.Action)
	} else if resp.Status == statusFailed {
		ob.Logger.Errorf("动作请求 `%v` 处理失败, 错误: %v", r.Action, resp.Message)
	} else {
		err := fmt.Errorf("动作处理器没有正确设置响应状态")
		ob.Logger.Warn(err)
		w.WriteFailed(RetCodeBadHandler, err)
	}
	return
}

func (ob *OneBot) decodeAndHandleRequest(actionBytes []byte, isBinary bool, comm RequestComm) Response {
	request, err := decodeRequest(actionBytes, isBinary, comm)
	if err != nil {
		err := fmt.Errorf("动作请求解析失败, 错误: %v", err)
		ob.Logger.Warn(err)
		return failedResponse(RetCodeBadRequest, err)
	}
	return ob.handleRequest(&request)
}

func (ob *OneBot) encodeResponse(resp Response, isBinary bool) ([]byte, error) {
	respBytes, err := resp.encode(isBinary)
	if err != nil {
		err := fmt.Errorf("动作响应编码失败, 错误: %v", err)
		ob.Logger.Warn(err)
		respBytes, _ = failedResponse(RetCodeBadHandler, err).encode(isBinary)
	}
	return respBytes, err
}
