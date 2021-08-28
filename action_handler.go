package libonebot

type Handler interface {
	HandleAction(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)
