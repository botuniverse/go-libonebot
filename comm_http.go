package libonebot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
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
	log.Debugf("HTTP request: %v", r)

	// reject unsupported methods
	if r.Method != "POST" && r.Method != "GET" {
		log.Warnf("Action 只支持通过 POST 方式请求")
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
		httpFail(w, RetCodeInvalidRequest, "Action 请求体 MIME 类型必须是 application/json")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		httpFail(w, RetCodeInvalidRequest, "Action 请求体读取失败: %v", err)
		return
	}

	response := comm.ob.handleAction(bytesToString(bodyBytes))
	json.NewEncoder(w).Encode(response)
}

func httpFail(w http.ResponseWriter, retcode int, errFormat string, args ...interface{}) {
	err := fmt.Errorf(errFormat, args...)
	log.Warn(err)
	json.NewEncoder(w).Encode(failedResponse(retcode, err))
}

func commStartHTTP(c ConfigCommHTTP, ob *OneBot) commCloser {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	log.Infof("正在启动 HTTP (%v)...", addr)

	comm := &httpComm{ob: ob}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("HTTP (%v) 启动失败, 错误: %v", addr, err)
		} else {
			log.Infof("HTTP (%v) 已关闭", addr)
		}
	}()

	return func() {
		if err := server.Shutdown(context.TODO() /* TODO */); err != nil {
			log.Errorf("HTTP (%v) 关闭失败, 错误: %v", addr, err)
		}
	}
}
