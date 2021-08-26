package comm

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/botuniverse/go-libonebot/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type httpComm struct {
}

func (comm *httpComm) handle(w http.ResponseWriter, r *http.Request) {
	log.Debugf("HTTP request: %v", r)

	// Reject unsupported methods
	if r.Method != "POST" && r.Method != "GET" {
		log.Warnf("Action 只支持通过 POST 方式请求")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Handle GET requests
	if r.Method == "GET" {
		w.Write([]byte("<h1>It works!</h1>"))
		return
	}

	// Reject unsupported content types
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		log.Warnf("Action 请求体 MIME 类型必须是 application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warnf("获取 Action 请求体失败, 错误: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := utils.BytesToString(bodyBytes)
	log.Debugf("HTTP request body: %v", body)
	if !gjson.Valid(body) {
		log.Warnf("Action 请求体不是合法的 JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	actionRequest := gjson.Parse(body)
	actionResponse := handleAction(actionRequest)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(utils.StringToBytes(actionResponse.String()))
}

// Start an HTTP communication task.
func StartHTTPTask(host string, port uint16) {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("正在启动 HTTP (%v)...", addr)

	comm := &httpComm{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			log.Errorf("HTTP (%v) 启动失败, 错误: %v", addr, err)
		}
		log.Infof("HTTP (%v) 已关闭", addr)
	}()
}
