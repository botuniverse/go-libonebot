package libonebot

import "fmt"

// ActionMux 将动作请求按动作名称分发到不同的 Handler 对象处理.
type ActionMux struct {
	handlers         map[string]Handler
	extendedHandlers map[string]Handler
}

// NewActionMux 创建一个新的 ActionMux 对象.
func NewActionMux() *ActionMux {
	return &ActionMux{
		handlers:         make(map[string]Handler),
		extendedHandlers: make(map[string]Handler),
	}
}

// HandleAction 为 ActionMux 实现 Handler 接口.
func (mux *ActionMux) HandleAction(w ResponseWriter, r *Request) {
	// return "ok" if otherwise explicitly set to "failed"
	w.WriteOK()

	var handlers *map[string]Handler
	if r.Action.IsExtended {
		handlers = &mux.extendedHandlers
	} else {
		handlers = &mux.handlers
	}

	handler := (*handlers)[r.Action.Name]
	if handler == nil {
		err := fmt.Errorf("动作 `%v` 不存在", r.Action)
		w.WriteFailed(RetCodeActionNotFound, err)
		return
	}

	handler.HandleAction(w, r)
}

// HandleFunc 将一个函数注册为指定核心动作的请求处理器.
func (mux *ActionMux) HandleFunc(action CoreAction, handler func(ResponseWriter, *Request)) {
	mux.Handle(action, HandlerFunc(handler))
}

// Handle 将一个 Handler 对象注册为指定核心动作的请求处理器.
func (mux *ActionMux) Handle(action CoreAction, handler Handler) {
	if action.name == "" {
		panic("动作名称不能为空")
	}
	mux.handlers[action.name] = handler
}

// HandleFuncExtended 将一个函数注册为指定扩展动作的请求处理器.
func (mux *ActionMux) HandleFuncExtended(action string, handler func(ResponseWriter, *Request)) {
	mux.HandleExtended(action, HandlerFunc(handler))
}

// HandleExtended 将一个 Handler 对象注册为指定扩展动作的请求处理器.
func (mux *ActionMux) HandleExtended(action string, handler HandlerFunc) {
	// if the prefix is empty, then the action name starts with "_"
	mux.extendedHandlers[action] = handler
}
