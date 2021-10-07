package libonebot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/tevino/abool/v2"
)

type httpComm struct {
	ob               *OneBot
	accessToken      string
	eventEnabled     bool
	eventBufferSize  uint32
	latestEvents     []marshaledEvent
	latestEventsLock *sync.Mutex
	latestEventsCond *sync.Cond
}

func (comm *httpComm) handle(w http.ResponseWriter, r *http.Request) {
	comm.ob.Logger.Debugf("HTTP request: %v", r)

	// reject unsupported methods
	if r.Method != "POST" {
		comm.ob.Logger.Errorf("动作请求不支持通过 %v 方式请求", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// authorization
	if comm.accessToken != "" {
		if r.Header.Get("Authorization") != "Bearer "+comm.accessToken {
			comm.ob.Logger.Errorf("动作请求头中的 Authorization 不匹配")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	var isBinary bool
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/json") {
		isBinary = false
		contentType = "application/json"
	} else if strings.HasPrefix(contentType, "application/msgpack") {
		isBinary = true
		contentType = "application/msgpack"
	} else {
		// reject unsupported content types
		comm.ob.Logger.Errorf("动作请求体 MIME 类型不支持")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// once we got the action HTTP request, we respond "200 OK"
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		comm.fail(w, RetCodeBadRequest, "动作请求体读取失败, 错误: %v", err)
		return
	}

	request, err := decodeRequest(bodyBytes, isBinary)
	if err != nil {
		comm.fail(w, RetCodeBadRequest, "动作请求体解析失败, 错误: %v", err)
		return
	}

	var response Response
	if comm.eventEnabled && request.Action == ActionGetLatestEvents {
		// special action: get_latest_events
		response = comm.handleGetLatestEvents(&request)
	} else {
		response = comm.ob.handleRequest(&request)
	}

	respBytes, _ := comm.ob.encodeResponse(response, isBinary)
	w.Write(respBytes)
}

func (comm *httpComm) handleGetLatestEvents(r *Request) (resp Response) {
	resp.Echo = r.Echo
	w := ResponseWriter{resp: &resp}

	timeout, err := r.Params.GetInt64("timeout")
	if err != nil {
		timeout = 0 // 0 for no wait
	}
	if timeout < 0 {
		w.WriteFailed(RetCodeBadParam, errors.New("`timeout` 参数值无效"))
		return
	}

	limit, err := r.Params.GetInt64("limit")
	if err != nil {
		limit = 0 // 0 for no limit
	}
	if limit < 0 {
		w.WriteFailed(RetCodeBadParam, errors.New("`limit` 参数值无效"))
		return
	}

	comm.latestEventsLock.Lock()
	defer comm.latestEventsLock.Unlock()

	if timeout > 0 && len(comm.latestEvents) == 0 {
		// wait for new events or timeout
		isTimeout := abool.New()
		timer := time.AfterFunc(time.Duration(timeout)*time.Second, func() {
			isTimeout.Set()
			comm.latestEventsCond.Broadcast() // wake up everyone because everyone may be out of time
			// but note, calling get_latest_events concurrently is undefined behavior
		})
		for {
			comm.latestEventsCond.Wait()
			if len(comm.latestEvents) > 0 || isTimeout.IsSet() {
				break
			}
		}
		timer.Stop()
	}

	eventCount := int64(len(comm.latestEvents))
	if limit == 0 || limit > eventCount {
		// if no limit, return all events
		limit = eventCount
	}
	events := make([]AnyEvent, 0)
	for _, event := range comm.latestEvents[:limit] {
		events = append(events, event.raw)
	}
	comm.latestEvents = comm.latestEvents[limit:]
	w.WriteData(events)
	return
}

func (comm *httpComm) fail(w http.ResponseWriter, retcode int, errFormat string, args ...interface{}) {
	err := fmt.Errorf(errFormat, args...)
	comm.ob.Logger.Warn(err)
	json.NewEncoder(w).Encode(failedResponse(retcode, err))
}

func commRunHTTP(c ConfigCommHTTP, ob *OneBot, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	ob.Logger.Infof("正在启动 HTTP (%v)...", addr)

	comm := &httpComm{
		ob:               ob,
		accessToken:      c.AccessToken,
		eventEnabled:     c.EventEnabled,
		eventBufferSize:  c.EventBufferSize,
		latestEvents:     make([]marshaledEvent, 0),
		latestEventsLock: &sync.Mutex{},
	}
	comm.latestEventsCond = sync.NewCond(comm.latestEventsLock)

	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ob.Logger.Errorf("HTTP (%v) 启动失败, 错误: %v", addr, err)
		}
	}()

	if comm.eventEnabled {
		eventChan := ob.openEventListenChan()
	loop:
		for {
			select {
			case event := <-eventChan:
				comm.latestEventsLock.Lock()
				if comm.eventBufferSize > 0 && len(comm.latestEvents) >= int(comm.eventBufferSize) {
					comm.latestEvents = append(comm.latestEvents[1:], event)
				} else {
					comm.latestEvents = append(comm.latestEvents, event)
				}
				comm.latestEventsLock.Unlock()
				comm.latestEventsCond.Signal() // notify someone to take the events
			case <-ctx.Done():
				ob.closeEventListenChan(eventChan)
				break loop
			}
		}
	} else {
		<-ctx.Done()
	}

	if err := server.Shutdown(context.TODO()); err != nil {
		ob.Logger.Errorf("HTTP (%v) 关闭失败, 错误: %v", addr, err)
	}
	ob.Logger.Infof("HTTP (%v) 已关闭", addr)
}
