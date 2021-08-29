package onebot

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type wsComm struct {
	onebot *OneBot
	addr   string
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (comm *wsComm) handle(w http.ResponseWriter, r *http.Request) {
	log.Infof("收到来自 %v 的 WebSocket (%v) 连接请求", r.RemoteAddr, comm.addr)
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("WebSocket (%v) 连接失败, 错误: %v", comm.addr, err)
		return
	}
	log.Infof("WebSocket (%v) 连接成功", comm.addr)
	defer conn.Close()

	// protect concurrent writes to the same connection
	connWriteLock := &sync.Mutex{}

	eventChan := comm.onebot.openEventListenChan()
	defer comm.onebot.closeEventListenChan(eventChan)
	go func() {
		// keep pushing events throught the connection
		for event := range eventChan {
			log.Debugf("通过 WebSocket (%v) 推送事件, %v", comm.addr, event.name)
			connWriteLock.Lock()
			conn.WriteMessage(websocket.TextMessage, event.bytes)
			connWriteLock.Unlock()
		}
	}()

	for {
		// this is the only one place we read from the connection, no need to lock
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Infof("WebSocket (%v) 连接断开", comm.addr)
			} else {
				log.Errorf("WebSocket (%v) 连接异常断开, 错误: %v", comm.addr, err)
			}
			break
		}

		response := comm.onebot.handleAction(bytesToString(messageBytes))
		connWriteLock.Lock()
		conn.WriteJSON(response)
		connWriteLock.Unlock()
	}
}

func commStartWS(host string, port uint16, onebot *OneBot) commCloser {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("正在启动 WebSocket (%v)...", addr)

	comm := &wsComm{onebot: onebot, addr: addr}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("WebSocket (%v) 启动失败, 错误: %v", addr, err)
		} else {
			log.Infof("WebSocket (%v) 已关闭", addr)
		}
	}()

	return func() {
		if err := server.Shutdown(context.TODO() /* TODO */); err != nil {
			log.Errorf("WebSocket (%v) 关闭失败, 错误: %v", addr, err)
		}
	}
}
