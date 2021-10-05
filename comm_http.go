package libonebot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type httpComm struct {
	ob               *OneBot
	latestEvents     []marshaledEvent
	latestEventsLock sync.Mutex
}

func (comm *httpComm) handleGET(w http.ResponseWriter, r *http.Request) {
	// TODO
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>It works!</h1>"))
}

func (comm *httpComm) handle(w http.ResponseWriter, r *http.Request) {
	comm.ob.Logger.Debugf("HTTP request: %v", r)

	// reject unsupported methods
	if r.Method != "POST" && r.Method != "GET" {
		comm.ob.Logger.Warnf("动作请求只支持通过 POST 方式请求")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// handle GET requests
	if r.Method == "GET" {
		comm.handleGET(w, r)
		return
	}

	// once we got the action HTTP request, we respond "200 OK"
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// reject unsupported content types
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		comm.fail(w, RetCodeBadRequest, "动作请求体 MIME 类型必须是 application/json")
		return
	}
	// TODO: Content-Type: application/msgpack

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		comm.fail(w, RetCodeBadRequest, "动作请求体读取失败: %v", err)
		return
	}

	request, err := parseActionRequest(bodyBytes, false)
	if err != nil {
		comm.fail(w, RetCodeBadRequest, "动作请求解析失败, 错误: %v", err)
		return
	}
	var response Response
	if request.Action == ActionGetLatestEvents {
		// special action: get_latest_events
		response = comm.handleGetLatestEvents(&request)
	} else {
		response = comm.ob.handleActionRequest(&request)
	}
	json.NewEncoder(w).Encode(response)
}

func (comm *httpComm) handleGetLatestEvents(r *Request) (resp Response) {
	resp.Echo = r.Echo
	w := ResponseWriter{resp: &resp}
	events := make([]AnyEvent, 0)
	// TODO: use condvar to wait until there are events
	comm.latestEventsLock.Lock()
	for _, event := range comm.latestEvents {
		events = append(events, event.raw)
	}
	comm.latestEvents = make([]marshaledEvent, 0)
	comm.latestEventsLock.Unlock()
	w.WriteData(events)
	return
}

func (comm *httpComm) fail(w http.ResponseWriter, retcode int, errFormat string, args ...interface{}) {
	err := fmt.Errorf(errFormat, args...)
	comm.ob.Logger.Warn(err)
	json.NewEncoder(w).Encode(failedResponse(retcode, err))
}

func commStartHTTP(c ConfigCommHTTP, ob *OneBot) commCloser {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	ob.Logger.Infof("正在启动 HTTP (%v)...", addr)

	comm := &httpComm{
		ob:           ob,
		latestEvents: make([]marshaledEvent, 0),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	eventChan := ob.openEventListenChan()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for event := range eventChan {
			comm.latestEventsLock.Lock()
			comm.latestEvents = append(comm.latestEvents, event)
			comm.latestEventsLock.Unlock()
		}
	}()

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ob.Logger.Errorf("HTTP (%v) 启动失败, 错误: %v", addr, err)
		} else {
			ob.Logger.Infof("HTTP (%v) 已关闭", addr)
		}
	}()

	return func() {
		ob.closeEventListenChan(eventChan)
		wg.Wait()
		if err := server.Shutdown(context.TODO() /* TODO */); err != nil {
			ob.Logger.Errorf("HTTP (%v) 关闭失败, 错误: %v", addr, err)
		}
		// TODO: wg.Wait() 后再输出已关闭
	}
}
