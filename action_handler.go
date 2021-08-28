package onebot

type Handler interface {
	HandleAction(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (handler HandlerFunc) HandleAction(w ResponseWriter, r *Request) {
	handler(w, r)
}
