package libonebot

// Handler 是动作处理器需要实现的接口.
type Handler interface {
	HandleAction(ResponseWriter, *Request)
}

// HandlerFunc 表示一个实现 Handler 接口的函数.
type HandlerFunc func(ResponseWriter, *Request)

// HandleAction 为 HandlerFunc 实现 Handler 接口.
func (handler HandlerFunc) HandleAction(w ResponseWriter, r *Request) {
	handler(w, r)
}
