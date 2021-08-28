package libonebot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/botuniverse/go-libonebot/utils"
	log "github.com/sirupsen/logrus"
)

type httpComm struct {
	actionMux *ActionMux
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
		comm.fail(w, RetCodeInvalidRequest, "Action 请求体 MIME 类型必须是 application/json")
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		comm.fail(w, RetCodeInvalidRequest, "Action 请求体读取失败: %v", err)
		return
	}

	response := comm.actionMux.HandleAction(utils.BytesToString(bodyBytes))
	json.NewEncoder(w).Encode(response)
}

func (comm *httpComm) fail(w http.ResponseWriter, retcode int, errFormat string, args ...interface{}) {
	err := fmt.Errorf(errFormat, args...)
	log.Warn(err)
	json.NewEncoder(w).Encode(failedResponse(retcode, err))
}

func commStartHTTP(host string, port uint16, actionMux *ActionMux) {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("正在启动 HTTP (%v)...", addr)

	comm := &httpComm{
		actionMux: actionMux,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			log.Errorf("HTTP (%v) 启动失败, 错误: %v", addr, err)
		}
		log.Infof("HTTP (%v) 已关闭", addr)
	}()
}
