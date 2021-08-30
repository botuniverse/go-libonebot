package libonebot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type httpComm struct {
	ob *OneBot
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
		comm.fail(w, RetCodeInvalidRequest, "动作请求体 MIME 类型必须是 application/json")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		comm.fail(w, RetCodeInvalidRequest, "动作请求体读取失败: %v", err)
		return
	}

	response := comm.ob.handleAction(bytesToString(bodyBytes))
	json.NewEncoder(w).Encode(response)
}

func (comm *httpComm) fail(w http.ResponseWriter, retcode int, errFormat string, args ...interface{}) {
	err := fmt.Errorf(errFormat, args...)
	comm.ob.Logger.Warn(err)
	json.NewEncoder(w).Encode(failedResponse(retcode, err))
}

func commStartHTTP(c ConfigCommHTTP, ob *OneBot) commCloser {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	ob.Logger.Infof("正在启动 HTTP (%v)...", addr)

	comm := &httpComm{ob: ob}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ob.Logger.Errorf("HTTP (%v) 启动失败, 错误: %v", addr, err)
		} else {
			ob.Logger.Infof("HTTP (%v) 已关闭", addr)
		}
	}()

	return func() {
		if err := server.Shutdown(context.TODO() /* TODO */); err != nil {
			ob.Logger.Errorf("HTTP (%v) 关闭失败, 错误: %v", addr, err)
		}
		// TODO: wg.Wait() 后再输出已关闭
	}
}
