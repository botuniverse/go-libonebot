package libonebot

import (
	"fmt"
	"sort"
)

// ActionMux 将动作请求按动作名称分发到不同的 Handler 对象处理.
type ActionMux struct {
	handlers map[string]Handler
}

// NewActionMux 创建一个新的 ActionMux 对象.
func NewActionMux() *ActionMux {
	mux := &ActionMux{
		handlers: make(map[string]Handler),
	}
	mux.HandleFunc(ActionGetSupportedActions, mux.handleGetSupportedActions)
	return mux
}

func (mux *ActionMux) handleGetSupportedActions(w ResponseWriter, r *Request) {
	actions := make([]string, 0, len(mux.handlers))
	for action := range mux.handlers {
		actions = append(actions, action)
	}
	sort.Slice(actions, func(i, j int) bool {
		return actions[i] < actions[j]
	})
	w.WriteData(actions)
}

// HandleAction 为 ActionMux 实现 Handler 接口.
func (mux *ActionMux) HandleAction(w ResponseWriter, r *Request) {
	// return "ok" if otherwise explicitly set to "failed"
	w.WriteOK()

	handler := mux.handlers[r.Action]
	if handler == nil {
		err := fmt.Errorf("动作 `%v` 不存在", r.Action)
		w.WriteFailed(RetCodeUnsupportedAction, err)
		return
	}

	handler.HandleAction(w, r)
}

// HandleFunc 将一个函数注册为指定动作的请求处理器.
//
// 若要注册为核心动作的请求处理器, 建议使用 ActionXxx 常量作为动作名.
func (mux *ActionMux) HandleFunc(action string, handler func(ResponseWriter, *Request)) {
	mux.Handle(action, HandlerFunc(handler))
}

// Handle 将一个 Handler 对象注册为指定动作的请求处理器.
//
// 若要注册为核心动作的请求处理器, 建议使用 ActionXxx 常量作为动作名.
func (mux *ActionMux) Handle(action string, handler Handler) {
	if action == "" {
		panic("动作名称不能为空")
	}
	mux.handlers[action] = handler
}
