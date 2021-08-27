package comm

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/botuniverse/go-libonebot/action"
	"github.com/botuniverse/go-libonebot/event"
	"github.com/botuniverse/go-libonebot/utils"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type wsComm struct {
	eventDispatcher *event.EventDispatcher
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (comm *wsComm) handle(w http.ResponseWriter, r *http.Request) {
	log.Infof("收到来自 %v 的 WebSocket 连接请求", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("WebSocket 连接失败, 错误: %v", err)
		return
	}
	log.Infof("WebSocket 连接成功")
	defer conn.Close()

	// protect concurrent writes to the same connection
	connWriteLock := &sync.Mutex{}

	eventChan := comm.eventDispatcher.OpenOutChan()
	defer comm.eventDispatcher.CloseOutChan(eventChan)
	go func() {
		// keep pushing events throught the connection
		for eventBytes := range eventChan {
			connWriteLock.Lock()
			conn.WriteMessage(websocket.TextMessage, eventBytes)
			connWriteLock.Unlock()
		}
	}()

	for {
		// this is the only one place we read from the connection, no need to lock
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Infof("WebSocket 连接断开")
			} else {
				log.Errorf("WebSocket 连接异常断开, 错误: %v", err)
			}
			break
		}

		message := utils.BytesToString(messageBytes)
		log.Debugf("WebSocket message: %v", message)
		if !gjson.Valid(message) {
			log.Warnf("Action 请求体不是合法的 JSON, 已忽略")
			continue
		}
		actionRequest := gjson.Parse(message)
		actionResponse := action.HandleAction(actionRequest)
		connWriteLock.Lock()
		conn.WriteMessage(websocket.TextMessage, utils.StringToBytes(actionResponse.String()))
		connWriteLock.Unlock()
	}
}

// Start a WebSocket commmunication task.
func StartWSTask(host string, port uint16, eventDispatcher *event.EventDispatcher) {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("正在启动 WebSocket (%v)...", addr)

	comm := &wsComm{eventDispatcher: eventDispatcher}
	mux := http.NewServeMux()
	mux.HandleFunc("/", comm.handle)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			log.Errorf("WebSocket (%v) 启动失败, 错误: %v", addr, err)
		}
		log.Infof("WebSocket (%v) 已关闭", addr)
	}()
}
